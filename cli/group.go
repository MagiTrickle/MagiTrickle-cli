package cli

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"magitrickle-terminal/pkg/api/types"

	"github.com/spf13/cobra"
)

var groupCmd = &cobra.Command{
	Use:   "group",
	Short: "Manage groups (list, create, update, delete, etc.)",
	Long: `Allows listing existing groups, creating new groups, updating or 
removing them. Under the hood, this command calls /api/v1/groups endpoints.`,
}

var listgroupCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List existing groups",
	Long: `Fetches a list of groups (optionally with rules) from /api/v1/groups. 
Use --with-rules to include rule details in the response.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		withRules, _ := cmd.Flags().GetBool("with-rules")
		url := "/api/v1/groups"
		if withRules {
			url += "?with_rules=true"
		}

		resp, err := doUnixRequest(http.MethodGet, url, nil)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return parseAPIError(resp)
		}

		var groupsRes types.GroupsRes
		if err := json.NewDecoder(resp.Body).Decode(&groupsRes); err != nil {
			return fmt.Errorf("failed to decode GroupsRes: %w", err)
		}

		if groupsRes.Groups == nil || len(*groupsRes.Groups) == 0 {
			fmt.Println("No groups found.")
			return nil
		}

		fmt.Println("Groups:")
		for _, g := range *groupsRes.Groups {
			fmt.Printf(" - ID: %s\n   Name: %s\n   Interface: %s\n   Enabled: %v\n   Color: %s\n",
				g.ID.String(), g.Name, g.Interface, g.Enable, g.Color)

			if withRules && g.Rules != nil && len(*g.Rules) > 0 {
				fmt.Println("   Rules:")
				for _, r := range *g.Rules {
					fmt.Printf("     * %s (%s) => %s [enabled: %v]\n",
						r.Name, r.Type, r.Rule, r.Enable)
				}
			}
			fmt.Println()
		}
		return nil
	},
}

var createGroupCmd = &cobra.Command{
	Use:     "create",
	Aliases: []string{"add"},
	Short:   "Create a new group",
	Long: `Creates a new group by sending a POST request to /api/v1/groups. 
You can specify name, interface, enable/disable, color, etc. If you pass --save,
the config will be saved immediately on the server side.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		name, _ := cmd.Flags().GetString("name")
		iface, _ := cmd.Flags().GetString("interface")
		enable, _ := cmd.Flags().GetBool("enable")
		colorStr, _ := cmd.Flags().GetString("color")

		reqBody := types.GroupReq{
			Name:      name,
			Interface: iface,
			Color:     colorStr,
			Enable:    &enable,
		}

		resp, err := doUnixJSON(http.MethodPost, "/api/v1/groups", reqBody)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return parseAPIError(resp)
		}

		var groupRes types.GroupRes
		if err := json.NewDecoder(resp.Body).Decode(&groupRes); err != nil {
			return fmt.Errorf("failed to decode GroupRes: %w", err)
		}

		fmt.Println("Group created successfully")
		fmt.Printf(" ID: %s\n Name: %s\n Interface: %s\n Enabled: %v\n Color: %s\n",
			groupRes.ID.String(), groupRes.Name, groupRes.Interface, groupRes.Enable, groupRes.Color)
		return nil
	},
}

var updateGroupCmd = &cobra.Command{
	Use:     "update",
	Aliases: []string{"edit"},
	Short:   "Update an existing group",
	Long: `Updates an existing group by sending a PUT request to /api/v1/groups/{groupID}. 
You must specify the group ID and optionally new name, interface, color, etc. 
Example:
    magitrickle group update <GROUP_ID> --name=NewName --enable=false --save
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("please provide the group ID as an argument, e.g. 'magitrickle group update <GROUP_ID>'")
		}
		groupID := args[0]

		name, _ := cmd.Flags().GetString("name")
		iface, _ := cmd.Flags().GetString("interface")
		enable, _ := cmd.Flags().GetBool("enable")
		colorStr, _ := cmd.Flags().GetString("color")

		reqBody := types.GroupReq{
			Name:      name,
			Interface: iface,
			Color:     colorStr,
			Enable:    &enable,
		}

		saveFlag, _ := cmd.Flags().GetBool("save")
		var urlBuilder strings.Builder
		urlBuilder.WriteString("/api/v1/groups/")
		urlBuilder.WriteString(groupID)
		if saveFlag {
			urlBuilder.WriteString("?save=true")
		}

		resp, err := doUnixJSON(http.MethodPut, urlBuilder.String(), reqBody)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return parseAPIError(resp)
		}

		var groupRes types.GroupRes
		if err := json.NewDecoder(resp.Body).Decode(&groupRes); err != nil {
			return fmt.Errorf("failed to decode updated GroupRes: %w", err)
		}

		fmt.Println("Group updated successfully")
		fmt.Printf(" ID: %s\n Name: %s\n Interface: %s\n Enabled: %v\n Color: %s\n",
			groupRes.ID.String(), groupRes.Name, groupRes.Interface, groupRes.Enable, groupRes.Color)
		return nil
	},
}

var deleteGroupCmd = &cobra.Command{
	Use:     "delete",
	Aliases: []string{"rm"},
	Short:   "Delete an existing group",
	Long: `Removes a group by sending a DELETE request to /api/v1/groups/{groupID}. 
Usage:
    magitrickle group delete <GROUP_ID> --save
If --save is specified, the server will persist configuration changes.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("please provide the group ID to delete, e.g. 'magitrickle group delete <GROUP_ID>'")
		}
		groupID := args[0]

		saveFlag, _ := cmd.Flags().GetBool("save")
		var urlBuilder strings.Builder
		urlBuilder.WriteString("/api/v1/groups/")
		urlBuilder.WriteString(groupID)
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

		fmt.Println("Group deleted successfully")
		return nil
	},
}

func init() {
	groupCmd.AddCommand(listgroupCmd)
	groupCmd.AddCommand(createGroupCmd)
	groupCmd.AddCommand(updateGroupCmd)
	groupCmd.AddCommand(deleteGroupCmd)

	listgroupCmd.Flags().Bool("with-rules", false, "Include rules for each group")

	createGroupCmd.Flags().String("name", "NewGroup", "Group name")
	createGroupCmd.Flags().String("interface", "br0", "Network interface for the group")
	createGroupCmd.Flags().Bool("enable", true, "Enable the group upon creation")
	createGroupCmd.Flags().String("color", "#ffffff", "Color hex code for the group")

	updateGroupCmd.Flags().String("name", "", "New name for the group")
	updateGroupCmd.Flags().String("interface", "", "New interface for the group")
	updateGroupCmd.Flags().Bool("enable", true, "Enable/disable the group")
	updateGroupCmd.Flags().String("color", "", "Color hex code for the group")
	updateGroupCmd.Flags().Bool("save", false, "Save config changes (append ?save=true)")

	deleteGroupCmd.Flags().Bool("save", false, "Save config changes (append ?save=true)")
}
