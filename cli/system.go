package cli

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Ponywka/MagiTrickle/backend/pkg/api/types"

	"github.com/spf13/cobra"
)

// systemCmd is the parent command for system-related operations
var systemCmd = &cobra.Command{
	Use:   "system",
	Short: "Manage system hooks, interfaces, and configuration",
	Long: `Provides commands to interact with netfilterd hooks, 
list available interfaces, and save the current configuration.`,
}

var netfilterdCmd = &cobra.Command{
	Use:   "netfilterd",
	Short: "Trigger a netfilter.d event",
	Long: `Triggers a netfilter.d hook by sending a POST request to 
/api/v1/system/hooks/netfilterd. You can specify the hook type and the table 
(e.g., "filter", "nat").`,
	RunE: func(cmd *cobra.Command, args []string) error {
		hookType, _ := cmd.Flags().GetString("type")
		table, _ := cmd.Flags().GetString("table")

		reqData := types.NetfilterDHookReq{
			Type:  hookType,
			Table: table,
		}
		resp, err := doUnixJSON(http.MethodPost, "/api/v1/system/hooks/netfilterd", reqData)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return parseAPIError(resp)
		}
		fmt.Println("Netfilterd hook triggered successfully")
		return nil
	},
}

var interfacesCmd = &cobra.Command{
	Use:   "interfaces",
	Short: "List network interfaces",
	Long: `Lists all available interfaces recognized by MagiTrickle
by sending a GET request to /api/v1/system/interfaces.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		resp, err := doUnixRequest(http.MethodGet, "/api/v1/system/interfaces", nil)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return parseAPIError(resp)
		}

		var ifaces types.InterfacesRes
		if err := json.NewDecoder(resp.Body).Decode(&ifaces); err != nil {
			return fmt.Errorf("failed to decode InterfacesRes: %w", err)
		}

		if len(ifaces.Interfaces) == 0 {
			fmt.Println("No interfaces found.")
			return nil
		}

		fmt.Println("Available Interfaces:")
		for _, iface := range ifaces.Interfaces {
			fmt.Printf("  - %s\n", iface.ID)
		}
		return nil
	},
}

var saveConfigCmd = &cobra.Command{
	Use:   "save-config",
	Short: "Save the current configuration",
	Long:  `Saves the current MagiTrickle configuration to persistent storage.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		resp, err := doUnixRequest(http.MethodPost, "/api/v1/system/config/save", nil)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
			return parseAPIError(resp)
		}
		fmt.Println("Configuration saved successfully")
		return nil
	},
}

func init() {
	systemCmd.AddCommand(netfilterdCmd)
	systemCmd.AddCommand(interfacesCmd)
	systemCmd.AddCommand(saveConfigCmd)

	netfilterdCmd.Flags().String("type", "filter", "Hook type (e.g., 'filter', 'nat')")
	netfilterdCmd.Flags().String("table", "filter", "Netfilter table (e.g., 'filter', 'nat')")
}
