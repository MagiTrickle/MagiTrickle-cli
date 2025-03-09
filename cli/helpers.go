package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Ponywka/MagiTrickle/backend/pkg/api/types"
)

func parseAPIError(resp *http.Response) error {
	defer resp.Body.Close()

	var errRes types.ErrorRes
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("request failed with status code %d (and body read error: %v)", resp.StatusCode, err)
	}

	if json.Unmarshal(body, &errRes) == nil && errRes.Error != "" {
		return fmt.Errorf("api error %d: %s", resp.StatusCode, errRes.Error)
	}

	return fmt.Errorf("request failed with status code %d", resp.StatusCode)
}
