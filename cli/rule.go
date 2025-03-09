package cli

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/Ponywka/MagiTrickle/backend/pkg/api/types"

	"github.com/spf13/cobra"
)

// ruleCmd – корневая команда для управления правилами (rules).
var ruleCmd = &cobra.Command{
	Use:   "rule",
	Short: "Manage rules within groups",
	Long: `Allows listing, creating, updating, deleting rules via 
/api/v1/groups/{groupID}/rules/... endpoints.`,
}

// listRulesCmd – GET /api/v1/groups/{groupID}/rules
// Выводит все правила группы.
var listRulesCmd = &cobra.Command{
	Use:     "list <GROUP_ID>",
	Aliases: []string{"ls"},
	Short:   "List all rules in the specified group",
	Long: `Calls GET /api/v1/groups/{groupID}/rules to retrieve all rules 
for the given group ID.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		groupID := args[0]
		url := fmt.Sprintf("/api/v1/groups/%s/rules", groupID)

		resp, err := doUnixRequest(http.MethodGet, url, nil)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return parseAPIError(resp)
		}

		var rulesRes types.RulesRes
		if err := json.NewDecoder(resp.Body).Decode(&rulesRes); err != nil {
			return fmt.Errorf("failed to decode RulesRes: %w", err)
		}

		if rulesRes.Rules == nil || len(*rulesRes.Rules) == 0 {
			fmt.Println("No rules found for this group.")
			return nil
		}

		fmt.Println("Rules in group", groupID, ":")
		for _, r := range *rulesRes.Rules {
			fmt.Printf(" - ID: %s | Name: %s | Type: %s | Rule: %s | Enabled: %v\n",
				r.ID.String(), r.Name, r.Type, r.Rule, r.Enable)
		}
		return nil
	},
}

// replaceRulesCmd – PUT /api/v1/groups/{groupID}/rules
// Полностью заменяет массив правил группы.
var replaceRulesCmd = &cobra.Command{
	Use:   "replace <GROUP_ID>",
	Short: "Replace all rules in a group with a new set",
	Long: `Calls PUT /api/v1/groups/{groupID}/rules to replace all rules in 
the given group. The new set of rules must be provided in a JSON file (via --file). 
If --save is used, changes will be persisted to config immediately.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		groupID := args[0]

		filePath, _ := cmd.Flags().GetString("file")
		if filePath == "" {
			return errors.New("please specify --file=<path_to_json> with an array of rules")
		}
		saveFlag, _ := cmd.Flags().GetBool("save")

		// Читаем содержимое JSON-файла
		content, err := os.ReadFile(filePath)
		if err != nil {
			return fmt.Errorf("failed to read file %s: %w", filePath, err)
		}

		// Парсим файл в types.RulesReq
		var rulesReq types.RulesReq
		if err := json.Unmarshal(content, &rulesReq); err != nil {
			return fmt.Errorf("failed to parse JSON from file: %w", err)
		}

		// Формируем URL
		var urlBuilder strings.Builder
		urlBuilder.WriteString("/api/v1/groups/")
		urlBuilder.WriteString(groupID)
		urlBuilder.WriteString("/rules")
		if saveFlag {
			urlBuilder.WriteString("?save=true")
		}

		resp, err := doUnixJSON(http.MethodPut, urlBuilder.String(), rulesReq)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return parseAPIError(resp)
		}

		var updated types.RulesRes
		if err := json.NewDecoder(resp.Body).Decode(&updated); err != nil {
			return fmt.Errorf("failed to decode updated RulesRes: %w", err)
		}

		fmt.Println("Rules replaced successfully. Current rules:")
		for _, r := range *updated.Rules {
			fmt.Printf(" - ID: %s | Name: %s | Type: %s | Rule: %s | Enabled: %v\n",
				r.ID.String(), r.Name, r.Type, r.Rule, r.Enable)
		}
		return nil
	},
}

// createRuleCmd – POST /api/v1/groups/{groupID}/rules
// Создаёт одно правило в группе.
var createRuleCmd = &cobra.Command{
	Use:   "create <GROUP_ID>",
	Short: "Create a single rule in the specified group",
	Long: `Calls POST /api/v1/groups/{groupID}/rules to create a single rule. 
Flags: 
  --name, --type, --rule, --enable, and optional --save.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		groupID := args[0]

		name, _ := cmd.Flags().GetString("name")
		rtype, _ := cmd.Flags().GetString("type")
		ruleStr, _ := cmd.Flags().GetString("rule")
		enable, _ := cmd.Flags().GetBool("enable")
		saveFlag, _ := cmd.Flags().GetBool("save")

		var urlBuilder strings.Builder
		urlBuilder.WriteString("/api/v1/groups/")
		urlBuilder.WriteString(groupID)
		urlBuilder.WriteString("/rules")
		if saveFlag {
			urlBuilder.WriteString("?save=true")
		}

		reqBody := types.RuleReq{
			Name:   name,
			Type:   rtype,
			Rule:   ruleStr,
			Enable: enable,
		}

		resp, err := doUnixJSON(http.MethodPost, urlBuilder.String(), reqBody)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return parseAPIError(resp)
		}

		var created types.RuleRes
		if err := json.NewDecoder(resp.Body).Decode(&created); err != nil {
			return fmt.Errorf("failed to decode created RuleRes: %w", err)
		}

		fmt.Println("Rule created successfully:")
		fmt.Printf(" ID: %s | Name: %s | Type: %s | Rule: %s | Enabled: %v\n",
			created.ID.String(), created.Name, created.Type, created.Rule, created.Enable)
		return nil
	},
}

// getRuleCmd – GET /api/v1/groups/{groupID}/rules/{ruleID}
// Возвращает конкретное правило из группы.
var getRuleCmd = &cobra.Command{
	Use:   "get <GROUP_ID> <RULE_ID>",
	Short: "Retrieve a specific rule by ID",
	Long: `Calls GET /api/v1/groups/{groupID}/rules/{ruleID} to fetch a single rule. 
You must provide both groupID and ruleID as arguments.`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		groupID := args[0]
		ruleID := args[1]

		url := fmt.Sprintf("/api/v1/groups/%s/rules/%s", groupID, ruleID)
		resp, err := doUnixRequest(http.MethodGet, url, nil)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return parseAPIError(resp)
		}

		var rule types.RuleRes
		if err := json.NewDecoder(resp.Body).Decode(&rule); err != nil {
			return fmt.Errorf("failed to decode RuleRes: %w", err)
		}

		fmt.Println("Rule info:")
		fmt.Printf(" ID: %s | Name: %s | Type: %s | Rule: %s | Enabled: %v\n",
			rule.ID.String(), rule.Name, rule.Type, rule.Rule, rule.Enable)
		return nil
	},
}

// updateRuleCmd – PUT /api/v1/groups/{groupID}/rules/{ruleID}
// Обновляет конкретное правило.
var updateRuleCmd = &cobra.Command{
	Use:   "update <GROUP_ID> <RULE_ID>",
	Short: "Update a specific rule by ID",
	Long: `Calls PUT /api/v1/groups/{groupID}/rules/{ruleID} to update a single rule. 
You can provide new values via flags (name, type, rule, enable). If --save is used, 
the config is persisted immediately.`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		groupID := args[0]
		ruleID := args[1]

		name, _ := cmd.Flags().GetString("name")
		rtype, _ := cmd.Flags().GetString("type")
		ruleStr, _ := cmd.Flags().GetString("rule")
		enable, _ := cmd.Flags().GetBool("enable")
		saveFlag, _ := cmd.Flags().GetBool("save")

		var urlBuilder strings.Builder
		urlBuilder.WriteString("/api/v1/groups/")
		urlBuilder.WriteString(groupID)
		urlBuilder.WriteString("/rules/")
		urlBuilder.WriteString(ruleID)
		if saveFlag {
			urlBuilder.WriteString("?save=true")
		}

		reqBody := types.RuleReq{
			Name:   name,
			Type:   rtype,
			Rule:   ruleStr,
			Enable: enable,
		}

		resp, err := doUnixJSON(http.MethodPut, urlBuilder.String(), reqBody)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return parseAPIError(resp)
		}

		var updated types.RuleRes
		if err := json.NewDecoder(resp.Body).Decode(&updated); err != nil {
			return fmt.Errorf("failed to decode updated RuleRes: %w", err)
		}

		fmt.Println("Rule updated successfully:")
		fmt.Printf(" ID: %s | Name: %s | Type: %s | Rule: %s | Enabled: %v\n",
			updated.ID.String(), updated.Name, updated.Type, updated.Rule, updated.Enable)
		return nil
	},
}

// deleteRuleCmd – DELETE /api/v1/groups/{groupID}/rules/{ruleID}
// Удаляет одно правило.
var deleteRuleCmd = &cobra.Command{
	Use:   "delete <GROUP_ID> <RULE_ID>",
	Short: "Delete a specific rule by ID",
	Long: `Calls DELETE /api/v1/groups/{groupID}/rules/{ruleID} to remove a single rule.
If --save is specified, configuration changes will be persisted immediately.`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		groupID := args[0]
		ruleID := args[1]

		saveFlag, _ := cmd.Flags().GetBool("save")
		var urlBuilder strings.Builder
		urlBuilder.WriteString("/api/v1/groups/")
		urlBuilder.WriteString(groupID)
		urlBuilder.WriteString("/rules/")
		urlBuilder.WriteString(ruleID)
		if saveFlag {
			urlBuilder.WriteString("?save=true")
		}

		resp, err := doUnixRequest(http.MethodDelete, urlBuilder.String(), nil)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return parseAPIError(resp)
		}

		fmt.Println("Rule deleted successfully")
		return nil
	},
}

func init() {
	// Регистрируем подкоманды у ruleCmd
	ruleCmd.AddCommand(listRulesCmd)
	ruleCmd.AddCommand(replaceRulesCmd)
	ruleCmd.AddCommand(createRuleCmd)
	ruleCmd.AddCommand(getRuleCmd)
	ruleCmd.AddCommand(updateRuleCmd)
	ruleCmd.AddCommand(deleteRuleCmd)

	// Флаги для "replace" (PUT /api/v1/groups/{groupID}/rules)
	// Ожидаем JSON-файл c массивом rules (types.RulesReq) через --file
	replaceRulesCmd.Flags().String("file", "", "Path to JSON file with an array of rules")
	replaceRulesCmd.Flags().Bool("save", false, "Save config changes (append ?save=true)")

	// Флаги для "create" (POST /api/v1/groups/{groupID}/rules)
	createRuleCmd.Flags().String("name", "", "Rule name")
	createRuleCmd.Flags().String("type", "domain", "Rule type (e.g. domain/ip/regex/etc.)")
	createRuleCmd.Flags().String("rule", "", "Rule value (e.g. example.com)")
	createRuleCmd.Flags().Bool("enable", true, "Enable this rule")
	createRuleCmd.Flags().Bool("save", false, "Save config changes (append ?save=true)")

	// Флаги для "update" (PUT /api/v1/groups/{groupID}/rules/{ruleID})
	updateRuleCmd.Flags().String("name", "", "New rule name")
	updateRuleCmd.Flags().String("type", "", "New rule type")
	updateRuleCmd.Flags().String("rule", "", "New rule value")
	updateRuleCmd.Flags().Bool("enable", true, "Enable/disable the rule")
	updateRuleCmd.Flags().Bool("save", false, "Save config changes (append ?save=true)")

	// Флаги для "delete" (DELETE /api/v1/groups/{groupID}/rules/{ruleID})
	deleteRuleCmd.Flags().Bool("save", false, "Save config changes (append ?save=true)")
}
