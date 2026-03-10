package wishlistsgateway

import (
	"bytes"
	"encoding/json"
	"fmt"
	"httpClient/internal/config"
	"httpClient/internal/http/handlers"
	"log/slog"
	"net/http"
	"time"
)

type WishlistGatewayClient struct {
	log        *slog.Logger
	cfg        *config.WishlistConfig
	httpClient *http.Client
}

func New(
	log *slog.Logger,
	fullCfg *config.Config,

) *WishlistGatewayClient {
	return &WishlistGatewayClient{
		log: log,
		cfg: &fullCfg.WishlistConfig,
		httpClient: &http.Client{
			Timeout: fullCfg.ServerConfig.Timeout * time.Second,
		},
	}
}

// Auth api calls
func (c *WishlistGatewayClient) Register(req handlers.RegisterRequestDTO) (handlers.AuthDTO, error) {
	const logPrefix = "wishlist.gateway.service.Register"

	baseURL := "http://" + c.cfg.GatewayAddres + ":" + c.cfg.GatewayPort + "/" + c.cfg.ApiURL
	url := baseURL + "/auth/register"

	log := c.log.With(
		slog.String("where", logPrefix),
		slog.String("url", url),
	)

	log.Info("Api Call Register")

	body, err := json.Marshal(req)
	if err != nil {
		log.Error("failed to marshal request", slog.String("err", err.Error()))
		return handlers.AuthDTO{}, fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := c.httpClient.Post(url, "application/json", bytes.NewReader(body))
	if err != nil {
		log.Error("failed to send request", slog.String("err", err.Error()))
		return handlers.AuthDTO{}, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		var errResp handlers.ErrorDTO
		if err := json.NewDecoder(resp.Body).Decode(&errResp); err != nil {
			log.Error("auth service error", slog.Int("status", resp.StatusCode))
			return handlers.AuthDTO{}, fmt.Errorf("auth service returned status %d", resp.StatusCode)
		}
		return handlers.AuthDTO{}, fmt.Errorf("wl service err: %s", errResp.Error)
	}

	var authResp handlers.AuthDTO
	if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
		log.Error("failed to decode response", slog.String("err", err.Error()))
		return handlers.AuthDTO{}, fmt.Errorf("failed to decode response: %w", err)
	}

	return authResp, nil
}

func (c *WishlistGatewayClient) Login(req handlers.LoginRequestDTO) (handlers.AuthDTO, error) {
	const logPrefix = "wishlist.gateway.service.Login"

	baseURL := "http://" + c.cfg.GatewayAddres + ":" + c.cfg.GatewayPort + "/" + c.cfg.ApiURL
	url := baseURL + "/auth/login"

	log := c.log.With(
		slog.String("where", logPrefix),
		slog.String("url", url),
	)

	log.Info("Api Call Register")

	body, err := json.Marshal(req)
	if err != nil {
		log.Error("failed to marshal request", slog.String("err", err.Error()))
		return handlers.AuthDTO{}, fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := c.httpClient.Post(url, "application/json", bytes.NewReader(body))
	if err != nil {
		log.Error("failed to send request", slog.String("err", err.Error()))
		return handlers.AuthDTO{}, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errResp handlers.ErrorDTO
		if err := json.NewDecoder(resp.Body).Decode(&errResp); err != nil {
			log.Error("auth service error", slog.Int("status", resp.StatusCode))
			return handlers.AuthDTO{}, fmt.Errorf("auth service returned status %d", resp.StatusCode)
		}
		return handlers.AuthDTO{}, fmt.Errorf("wl service err: %s", errResp.Error)
	}

	var authResp handlers.AuthDTO
	if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
		log.Error("failed to decode response", slog.String("err", err.Error()))
		return handlers.AuthDTO{}, fmt.Errorf("failed to decode response: %w", err)
	}

	return authResp, nil
}

func (c *WishlistGatewayClient) GetCurrentUser(token string) (handlers.UserInfoDTO, error) {
	const logPrefix = "wishlist.gateway.service.GetCurrentUser"

	baseURL := "http://" + c.cfg.GatewayAddres + ":" + c.cfg.GatewayPort + "/" + c.cfg.ApiURL
	url := baseURL + "/auth/self"

	log := c.log.With(
		slog.String("where", logPrefix),
		slog.String("url", url),
	)

	log.Info("Api Call GetCurrentUser")

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Error("failed to create request", slog.String("err", err.Error()))
		return handlers.UserInfoDTO{}, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		log.Error("failed to send request", slog.String("err", err.Error()))
		return handlers.UserInfoDTO{}, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errResp handlers.ErrorDTO
		if err := json.NewDecoder(resp.Body).Decode(&errResp); err != nil {
			log.Error("auth service error", slog.Int("status", resp.StatusCode))
			return handlers.UserInfoDTO{}, fmt.Errorf("auth service returned status %d", resp.StatusCode)
		}
		return handlers.UserInfoDTO{}, fmt.Errorf("wl service err: %s", errResp.Error)
	}

	var userInfo handlers.UserInfoDTO
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		log.Error("failed to decode response", slog.String("err", err.Error()))
		return handlers.UserInfoDTO{}, fmt.Errorf("failed to decode response: %w", err)
	}

	return userInfo, nil
}

func (c *WishlistGatewayClient) GetUserInfo(token string, userID string) (handlers.UserInfoDTO, error) {
	const logPrefix = "wishlist.gateway.service.GetUserInfo"

	baseURL := "http://" + c.cfg.GatewayAddres + ":" + c.cfg.GatewayPort + "/" + c.cfg.ApiURL
	//url := baseURL + "/auth/user"
	url := baseURL + fmt.Sprintf("/auth/user/%s", userID)

	log := c.log.With(
		slog.String("where", logPrefix),
		slog.String("url", url),
	)

	log.Info("Api Call GetUserInfo")

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Error("failed to create request", slog.String("err", err.Error()))
		return handlers.UserInfoDTO{}, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		log.Error("failed to send request", slog.String("err", err.Error()))
		return handlers.UserInfoDTO{}, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errResp handlers.ErrorDTO
		if err := json.NewDecoder(resp.Body).Decode(&errResp); err != nil {
			log.Error("auth service error", slog.Int("status", resp.StatusCode))
			return handlers.UserInfoDTO{}, fmt.Errorf("auth service returned status %d", resp.StatusCode)
		}
		return handlers.UserInfoDTO{}, fmt.Errorf("wl service err: %s", errResp.Error)
	}

	var userInfo handlers.UserInfoDTO
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		log.Error("failed to decode response", slog.String("err", err.Error()))
		return handlers.UserInfoDTO{}, fmt.Errorf("failed to decode response: %w", err)
	}

	return userInfo, nil
}

// wishlists api calls
func (c *WishlistGatewayClient) CreateWishlist(token string, req handlers.CreateWishlistRequestDTO) (handlers.WishlistDTO, error) {
	const logPrefix = "wishlist.gateway.service.CreateWishlist"

	baseURL := "http://" + c.cfg.GatewayAddres + ":" + c.cfg.GatewayPort + "/" + c.cfg.ApiURL
	url := baseURL + "/wishlists"

	log := c.log.With(
		slog.String("where", logPrefix),
		slog.String("url", url),
	)

	log.Info("Api Call CreateWishlist")

	body, err := json.Marshal(req)
	if err != nil {
		log.Error("failed to marshal request", slog.String("err", err.Error()))
		return handlers.WishlistDTO{}, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequest("POST", url, bytes.NewReader(body))
	if err != nil {
		log.Error("failed to create request", slog.String("err", err.Error()))
		return handlers.WishlistDTO{}, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+token)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		log.Error("failed to send request", slog.String("err", err.Error()))
		return handlers.WishlistDTO{}, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		var errResp handlers.ErrorDTO
		if err := json.NewDecoder(resp.Body).Decode(&errResp); err != nil {
			log.Error("wishlist service error", slog.Int("status", resp.StatusCode))
			return handlers.WishlistDTO{}, fmt.Errorf("wishlist service returned status %d", resp.StatusCode)
		}
		return handlers.WishlistDTO{}, fmt.Errorf("wl service err: %s", errResp.Error)
	}

	var wishlist handlers.WishlistDTO
	if err := json.NewDecoder(resp.Body).Decode(&wishlist); err != nil {
		log.Error("failed to decode response", slog.String("err", err.Error()))
		return handlers.WishlistDTO{}, fmt.Errorf("failed to decode response: %w", err)
	}

	return wishlist, nil
}

func (c *WishlistGatewayClient) GetUserWishlists(token, userID string, page, limit int32) (handlers.WishlistListDTO, error) {
	const logPrefix = "wishlist.gateway.service.GetUserWishlists"

	baseURL := "http://" + c.cfg.GatewayAddres + ":" + c.cfg.GatewayPort + "/" + c.cfg.ApiURL
	url := baseURL + fmt.Sprintf("/wishlists/user/%s?page=%d&page_size=%d", userID, page, limit)

	log := c.log.With(
		slog.String("where", logPrefix),
		slog.String("url", url),
	)

	log.Info("Api Call GetUserWishlists")

	httpReq, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Error("failed to create request", slog.String("err", err.Error()))
		return handlers.WishlistListDTO{}, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Authorization", "Bearer "+token)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		log.Error("failed to send request", slog.String("err", err.Error()))
		return handlers.WishlistListDTO{}, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errResp handlers.ErrorDTO
		if err := json.NewDecoder(resp.Body).Decode(&errResp); err != nil {
			log.Error("wishlist service error", slog.Int("status", resp.StatusCode))
			return handlers.WishlistListDTO{}, fmt.Errorf("wishlist service returned status %d", resp.StatusCode)
		}
		return handlers.WishlistListDTO{}, fmt.Errorf("wl service err: %s", errResp.Error)
	}

	var listResp handlers.WishlistListDTO
	if err := json.NewDecoder(resp.Body).Decode(&listResp); err != nil {
		log.Error("failed to decode response", slog.String("err", err.Error()))
		return handlers.WishlistListDTO{}, fmt.Errorf("failed to decode response: %w", err)
	}

	return listResp, nil
}

func (c *WishlistGatewayClient) GetPublicWishlists(token string, page, limit int32) (handlers.WishlistListDTO, error) {
	const logPrefix = "wishlist.gateway.service.GetPublicWishlists"

	baseURL := "http://" + c.cfg.GatewayAddres + ":" + c.cfg.GatewayPort + "/" + c.cfg.ApiURL
	url := baseURL + fmt.Sprintf("/wishlists/public?page=%d&page_size=%d", page, limit)

	log := c.log.With(
		slog.String("where", logPrefix),
		slog.String("url", url),
	)

	log.Info("Api Call GetPublicWishlists")

	httpReq, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Error("failed to create request", slog.String("err", err.Error()))
		return handlers.WishlistListDTO{}, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Authorization", "Bearer "+token)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		log.Error("failed to send request", slog.String("err", err.Error()))
		return handlers.WishlistListDTO{}, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errResp handlers.ErrorDTO
		if err := json.NewDecoder(resp.Body).Decode(&errResp); err != nil {
			log.Error("wishlist service error", slog.Int("status", resp.StatusCode))
			return handlers.WishlistListDTO{}, fmt.Errorf("wishlist service returned status %d", resp.StatusCode)
		}
		return handlers.WishlistListDTO{}, fmt.Errorf("wl service err: %s", errResp.Error)
	}

	var listResp handlers.WishlistListDTO
	if err := json.NewDecoder(resp.Body).Decode(&listResp); err != nil {
		log.Error("failed to decode response", slog.String("err", err.Error()))
		return handlers.WishlistListDTO{}, fmt.Errorf("failed to decode response: %w", err)
	}

	return listResp, nil
}

func (c *WishlistGatewayClient) GetWishlist(token, wishlistID string) (handlers.WishlistDTO, error) {
	const logPrefix = "wishlist.gateway.service.GetWishlist"

	baseURL := "http://" + c.cfg.GatewayAddres + ":" + c.cfg.GatewayPort + "/" + c.cfg.ApiURL
	url := baseURL + fmt.Sprintf("/wishlists/%s", wishlistID)

	log := c.log.With(
		slog.String("where", logPrefix),
		slog.String("url", url),
	)

	log.Info("Api Call GetWishlist")

	httpReq, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Error("failed to create request", slog.String("err", err.Error()))
		return handlers.WishlistDTO{}, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Authorization", "Bearer "+token)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		log.Error("failed to send request", slog.String("err", err.Error()))
		return handlers.WishlistDTO{}, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errResp handlers.ErrorDTO
		if err := json.NewDecoder(resp.Body).Decode(&errResp); err != nil {
			log.Error("wishlist service error", slog.Int("status", resp.StatusCode))
			return handlers.WishlistDTO{}, fmt.Errorf("wishlist service returned status %d", resp.StatusCode)
		}
		return handlers.WishlistDTO{}, fmt.Errorf("wl service err: %s", errResp.Error)
	}

	var wishlist handlers.WishlistDTO
	if err := json.NewDecoder(resp.Body).Decode(&wishlist); err != nil {
		log.Error("failed to decode response", slog.String("err", err.Error()))
		return handlers.WishlistDTO{}, fmt.Errorf("failed to decode response: %w", err)
	}

	return wishlist, nil
}

func (c *WishlistGatewayClient) UpdateWishlist(token, wishlistID string, req handlers.UpdateWishlistRequestDTO) (handlers.WishlistDTO, error) {
	const logPrefix = "wishlist.gateway.service.UpdateWishlist"

	baseURL := "http://" + c.cfg.GatewayAddres + ":" + c.cfg.GatewayPort + "/" + c.cfg.ApiURL
	url := baseURL + fmt.Sprintf("/wishlists/%s", wishlistID)

	log := c.log.With(
		slog.String("where", logPrefix),
		slog.String("url", url),
	)

	log.Info("Api Call UpdateWishlist")

	body, err := json.Marshal(req)
	if err != nil {
		log.Error("failed to marshal request", slog.String("err", err.Error()))
		return handlers.WishlistDTO{}, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequest("PUT", url, bytes.NewReader(body))
	if err != nil {
		log.Error("failed to create request", slog.String("err", err.Error()))
		return handlers.WishlistDTO{}, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+token)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		log.Error("failed to send request", slog.String("err", err.Error()))
		return handlers.WishlistDTO{}, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errResp handlers.ErrorDTO
		if err := json.NewDecoder(resp.Body).Decode(&errResp); err != nil {
			log.Error("wishlist service error", slog.Int("status", resp.StatusCode))
			return handlers.WishlistDTO{}, fmt.Errorf("wishlist service returned status %d", resp.StatusCode)
		}
		return handlers.WishlistDTO{}, fmt.Errorf("wl service err: %s", errResp.Error)
	}

	var wishlist handlers.WishlistDTO
	if err := json.NewDecoder(resp.Body).Decode(&wishlist); err != nil {
		log.Error("failed to decode response", slog.String("err", err.Error()))
		return handlers.WishlistDTO{}, fmt.Errorf("failed to decode response: %w", err)
	}

	return wishlist, nil
}

func (c *WishlistGatewayClient) DeleteWishlist(token, wishlistID string) error {
	const logPrefix = "wishlist.gateway.service.DeleteWishlist"

	baseURL := "http://" + c.cfg.GatewayAddres + ":" + c.cfg.GatewayPort + "/" + c.cfg.ApiURL
	url := baseURL + fmt.Sprintf("/wishlists/%s", wishlistID)

	log := c.log.With(
		slog.String("where", logPrefix),
		slog.String("url", url),
	)

	log.Info("Api Call DeleteWishlist")

	httpReq, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		log.Error("failed to create request", slog.String("err", err.Error()))
		return fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Authorization", "Bearer "+token)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		log.Error("failed to send request", slog.String("err", err.Error()))
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		var errResp handlers.ErrorDTO
		if err := json.NewDecoder(resp.Body).Decode(&errResp); err != nil {
			log.Error("wishlist service error", slog.Int("status", resp.StatusCode))
			return fmt.Errorf("wishlist service returned status %d", resp.StatusCode)
		}
		return fmt.Errorf("wl service err: %s", errResp.Error)
	}

	return nil
}

// item api calls
func (c *WishlistGatewayClient) CreateItem(token, wishlistID string, req handlers.CreateItemRequestDTO) (handlers.ItemDTO, error) {
	const logPrefix = "wishlist.gateway.service.CreateItem"

	baseURL := "http://" + c.cfg.GatewayAddres + ":" + c.cfg.GatewayPort + "/" + c.cfg.ApiURL
	url := baseURL + fmt.Sprintf("/wishlists/%s/items", wishlistID)

	log := c.log.With(
		slog.String("where", logPrefix),
		slog.String("url", url),
	)

	log.Info("Api Call CreateItem")

	body, err := json.Marshal(req)
	if err != nil {
		log.Error("failed to marshal request", slog.String("err", err.Error()))
		return handlers.ItemDTO{}, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequest("POST", url, bytes.NewReader(body))
	if err != nil {
		log.Error("failed to create request", slog.String("err", err.Error()))
		return handlers.ItemDTO{}, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+token)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		log.Error("failed to send request", slog.String("err", err.Error()))
		return handlers.ItemDTO{}, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		var errResp handlers.ErrorDTO
		if err := json.NewDecoder(resp.Body).Decode(&errResp); err != nil {
			log.Error("item service error", slog.Int("status", resp.StatusCode))
			return handlers.ItemDTO{}, fmt.Errorf("item service returned status %d", resp.StatusCode)
		}
		return handlers.ItemDTO{}, fmt.Errorf("wl service err: %s", errResp.Error)
	}

	var item handlers.ItemDTO
	if err := json.NewDecoder(resp.Body).Decode(&item); err != nil {
		log.Error("failed to decode response", slog.String("err", err.Error()))
		return handlers.ItemDTO{}, fmt.Errorf("failed to decode response: %w", err)
	}

	return item, nil
}

func (c *WishlistGatewayClient) GetItems(token, wishlistID string, page, limit int32) (handlers.ItemListDTO, error) {
	const logPrefix = "wishlist.gateway.service.GetItems"

	baseURL := "http://" + c.cfg.GatewayAddres + ":" + c.cfg.GatewayPort + "/" + c.cfg.ApiURL
	url := baseURL + fmt.Sprintf("/wishlists/%s/items?page=%d&page_size=%d", wishlistID, page, limit)

	log := c.log.With(
		slog.String("where", logPrefix),
		slog.String("url", url),
	)

	log.Info("Api Call GetItems")

	httpReq, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Error("failed to create request", slog.String("err", err.Error()))
		return handlers.ItemListDTO{}, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Authorization", "Bearer "+token)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		log.Error("failed to send request", slog.String("err", err.Error()))
		return handlers.ItemListDTO{}, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errResp handlers.ErrorDTO
		if err := json.NewDecoder(resp.Body).Decode(&errResp); err != nil {
			log.Error("item service error", slog.Int("status", resp.StatusCode))
			return handlers.ItemListDTO{}, fmt.Errorf("item service returned status %d", resp.StatusCode)
		}
		return handlers.ItemListDTO{}, fmt.Errorf("wl service err: %s", errResp.Error)
	}

	var listResp handlers.ItemListDTO
	if err := json.NewDecoder(resp.Body).Decode(&listResp); err != nil {
		log.Error("failed to decode response", slog.String("err", err.Error()))
		return handlers.ItemListDTO{}, fmt.Errorf("failed to decode response: %w", err)
	}

	return listResp, nil
}

func (c *WishlistGatewayClient) GetItem(token, itemID string) (handlers.ItemDTO, error) {
	const logPrefix = "wishlist.gateway.service.GetItem"

	baseURL := "http://" + c.cfg.GatewayAddres + ":" + c.cfg.GatewayPort + "/" + c.cfg.ApiURL
	url := baseURL + fmt.Sprintf("/items/%s", itemID)

	log := c.log.With(
		slog.String("where", logPrefix),
		slog.String("url", url),
	)

	log.Info("Api Call GetItem")

	httpReq, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Error("failed to create request", slog.String("err", err.Error()))
		return handlers.ItemDTO{}, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Authorization", "Bearer "+token)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		log.Error("failed to send request", slog.String("err", err.Error()))
		return handlers.ItemDTO{}, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errResp handlers.ErrorDTO
		if err := json.NewDecoder(resp.Body).Decode(&errResp); err != nil {
			log.Error("item service error", slog.Int("status", resp.StatusCode))
			return handlers.ItemDTO{}, fmt.Errorf("item service returned status %d", resp.StatusCode)
		}
		return handlers.ItemDTO{}, fmt.Errorf("wl service err: %s", errResp.Error)
	}

	var item handlers.ItemDTO
	if err := json.NewDecoder(resp.Body).Decode(&item); err != nil {
		log.Error("failed to decode response", slog.String("err", err.Error()))
		return handlers.ItemDTO{}, fmt.Errorf("failed to decode response: %w", err)
	}

	return item, nil
}

func (c *WishlistGatewayClient) UpdateItem(token, itemID string, req handlers.UpdateItemRequestDTO) (handlers.ItemDTO, error) {
	const logPrefix = "wishlist.gateway.service.UpdateItem"

	baseURL := "http://" + c.cfg.GatewayAddres + ":" + c.cfg.GatewayPort + "/" + c.cfg.ApiURL
	url := baseURL + fmt.Sprintf("/items/%s", itemID)

	log := c.log.With(
		slog.String("where", logPrefix),
		slog.String("url", url),
	)

	log.Info("Api Call UpdateItem")
	//url := fmt.Sprintf("%s/api/v1/items/%s", s.baseURL, itemID)

	body, err := json.Marshal(req)
	if err != nil {
		log.Error("failed to marshal request", slog.String("err", err.Error()))
		return handlers.ItemDTO{}, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequest("PUT", url, bytes.NewReader(body))
	if err != nil {
		log.Error("failed to create request", slog.String("err", err.Error()))
		return handlers.ItemDTO{}, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+token)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		log.Error("failed to send request", slog.String("err", err.Error()))
		return handlers.ItemDTO{}, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errResp handlers.ErrorDTO
		if err := json.NewDecoder(resp.Body).Decode(&errResp); err != nil {
			log.Error("item service error", slog.Int("status", resp.StatusCode))
			return handlers.ItemDTO{}, fmt.Errorf("item service returned status %d", resp.StatusCode)
		}
		return handlers.ItemDTO{}, fmt.Errorf("wl service err: %s", errResp.Error)
	}

	var item handlers.ItemDTO
	if err := json.NewDecoder(resp.Body).Decode(&item); err != nil {
		log.Error("failed to decode response", slog.String("err", err.Error()))
		return handlers.ItemDTO{}, fmt.Errorf("failed to decode response: %w", err)
	}

	return item, nil
}

func (c *WishlistGatewayClient) DeleteItem(token, itemID string) error {
	const logPrefix = "wishlist.gateway.service.DeleteItem"

	baseURL := "http://" + c.cfg.GatewayAddres + ":" + c.cfg.GatewayPort + "/" + c.cfg.ApiURL
	url := baseURL + fmt.Sprintf("/items/%s", itemID)

	log := c.log.With(
		slog.String("where", logPrefix),
		slog.String("url", url),
	)

	log.Info("Api Call DeleteItem")

	//url := fmt.Sprintf("%s/api/v1/items/%s", s.baseURL, itemID)

	httpReq, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		log.Error("failed to create request", slog.String("err", err.Error()))
		return fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Authorization", "Bearer "+token)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		log.Error("failed to send request", slog.String("err", err.Error()))
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		var errResp handlers.ErrorDTO
		if err := json.NewDecoder(resp.Body).Decode(&errResp); err != nil {
			log.Error("item service error", slog.Int("status", resp.StatusCode))
			return fmt.Errorf("item service returned status %d", resp.StatusCode)
		}
		return fmt.Errorf("wl service err: %s", errResp.Error)
	}

	return nil
}

func (c *WishlistGatewayClient) BookItem(token, itemID string) (handlers.ItemDTO, error) {
	const logPrefix = "wishlist.gateway.service.BookItem"

	baseURL := "http://" + c.cfg.GatewayAddres + ":" + c.cfg.GatewayPort + "/" + c.cfg.ApiURL
	url := baseURL + fmt.Sprintf("/items/%s/book", itemID)

	log := c.log.With(
		slog.String("where", logPrefix),
		slog.String("url", url),
	)

	log.Info("Api Call BookItem")
	//url := fmt.Sprintf("%s/api/v1/items/%s/book", s.baseURL, itemID)

	httpReq, err := http.NewRequest("POST", url, nil)
	if err != nil {
		log.Error("failed to create request", slog.String("err", err.Error()))
		return handlers.ItemDTO{}, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Authorization", "Bearer "+token)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		log.Error("failed to send request", slog.String("err", err.Error()))
		return handlers.ItemDTO{}, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errResp handlers.ErrorDTO
		if err := json.NewDecoder(resp.Body).Decode(&errResp); err != nil {
			log.Error("item service error", slog.Int("status", resp.StatusCode))
			return handlers.ItemDTO{}, fmt.Errorf("item service returned status %d", resp.StatusCode)
		}
		return handlers.ItemDTO{}, fmt.Errorf("wl service err: %s", errResp.Error)
	}

	var item handlers.ItemDTO
	if err := json.NewDecoder(resp.Body).Decode(&item); err != nil {
		log.Error("failed to decode response", slog.String("err", err.Error()))
		return handlers.ItemDTO{}, fmt.Errorf("failed to decode response: %w", err)
	}

	return item, nil
}

func (c *WishlistGatewayClient) UnbookItem(token, itemID string) (handlers.ItemDTO, error) {
	const logPrefix = "wishlist.gateway.service.UnbookItem"

	baseURL := "http://" + c.cfg.GatewayAddres + ":" + c.cfg.GatewayPort + "/" + c.cfg.ApiURL
	url := baseURL + fmt.Sprintf("/items/%s/unbook", itemID)

	log := c.log.With(
		slog.String("where", logPrefix),
		slog.String("url", url),
	)

	log.Info("Api Call UnbookItem")
	//url := fmt.Sprintf("%s/api/v1/items/%s/unbook", s.baseURL, itemID)

	httpReq, err := http.NewRequest("POST", url, nil)
	if err != nil {
		log.Error("failed to create request", slog.String("err", err.Error()))
		return handlers.ItemDTO{}, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Authorization", "Bearer "+token)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		log.Error("failed to send request", slog.String("err", err.Error()))
		return handlers.ItemDTO{}, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errResp handlers.ErrorDTO
		if err := json.NewDecoder(resp.Body).Decode(&errResp); err != nil {
			log.Error("item service error", slog.Int("status", resp.StatusCode))
			return handlers.ItemDTO{}, fmt.Errorf("item service returned status %d", resp.StatusCode)
		}
		return handlers.ItemDTO{}, fmt.Errorf("wl service err: %s", errResp.Error)
	}

	var item handlers.ItemDTO
	if err := json.NewDecoder(resp.Body).Decode(&item); err != nil {
		log.Error("failed to decode response", slog.String("err", err.Error()))
		return handlers.ItemDTO{}, fmt.Errorf("failed to decode response: %w", err)
	}

	return item, nil
}

// booking call api
func (c *WishlistGatewayClient) GetUserBookings(token string) (handlers.BookingListDTO, error) {
	const logPrefix = "wishlist.gateway.service.UnbookItem"

	baseURL := "http://" + c.cfg.GatewayAddres + ":" + c.cfg.GatewayPort + "/" + c.cfg.ApiURL
	url := baseURL + "/bookings"

	log := c.log.With(
		slog.String("where", logPrefix),
		slog.String("url", url),
	)

	log.Info("Api Call UnbookItem")

	httpReq, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Error("failed to create request", slog.String("err", err.Error()))
		return handlers.BookingListDTO{}, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Authorization", "Bearer "+token)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		log.Error("failed to send request", slog.String("err", err.Error()))
		return handlers.BookingListDTO{}, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errResp handlers.ErrorDTO
		if err := json.NewDecoder(resp.Body).Decode(&errResp); err != nil {
			log.Error("booking service error", slog.Int("status", resp.StatusCode))
			return handlers.BookingListDTO{}, fmt.Errorf("booking service returned status %d", resp.StatusCode)
		}
		return handlers.BookingListDTO{}, fmt.Errorf("wl service err: %s", errResp.Error)
	}

	var listResp handlers.BookingListDTO
	if err := json.NewDecoder(resp.Body).Decode(&listResp); err != nil {
		log.Error("failed to decode response", slog.String("err", err.Error()))
		return handlers.BookingListDTO{}, fmt.Errorf("failed to decode response: %w", err)
	}

	return listResp, nil
}
