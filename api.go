package hypixel

import (
	"bytes"
	"io"
	"net/http"
)

// Send Hypixel API HTTP Request
func (c *Client) Send(method string, head http.Header, path string, params *Params, payload ...byte) (*http.Response, error) {
	if method == "" {
		method = http.MethodGet
	}
	full := c.GetFullPath(path)
	if params != nil {
		full = params.String(full)
	}
	req, err := http.NewRequest(method, full,
		func() io.Reader {
			if payload != nil {
				return bytes.NewReader(payload)
			}
			return nil
		}(),
	)
	if err != nil {
		return nil, err
	}
	if head != nil {
		req.Header = head
	}
	if c.rate != nil {
		c.rate.WaitIfNeeded()
	}
	rsp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	if c.rate != nil {
		_ = c.rate.UpdateFromHeaders(rsp.Header)
	}
	return rsp, nil
}

// Authentication Add api key to header
//
// https://api.hypixel.net/#section/Authentication/ApiKey
func (c *Client) Authentication(header ...http.Header) http.Header {
	var h http.Header
	if len(header) == 0 {
		h = http.Header{}
	} else {
		h = header[0]
	}
	h.Set("API-Key", c.apiKey)
	return h
}

// GetPlayerData Data of a specific player, including game stats
// NEED API Key
//
// https://api.hypixel.net/#tag/Player-Data
func (c *Client) GetPlayerData(uuid string) (*http.Response, error) {
	return c.Send(http.MethodGet, c.Authentication(), "player", &Params{
		"uuid": uuid,
	})
}

// GetRecentGames The recently played games of a specific player
// NEED API Key
//
// https://api.hypixel.net/#tag/Player-Data/paths/~1v2~1recentgames/get
func (c *Client) GetRecentGames(uuid string) (*http.Response, error) {
	return c.Send(http.MethodGet, c.Authentication(), "recentgames", &Params{
		"uuid": uuid,
	})
}

// GetStatus The current online status of a specific player
// NEED API Key
//
// https://api.hypixel.net/#tag/Player-Data/paths/~1v2~1status/get
func (c *Client) GetStatus(uuid string) (*http.Response, error) {
	return c.Send(http.MethodGet, c.Authentication(), "status", &Params{
		"uuid": uuid,
	})
}

// GetGuild Retrieve a Guild by a player, id, or name
// NEED API Key
//
// https://api.hypixel.net/#tag/Player-Data/paths/~1v2~1guild/get
func (c *Client) GetGuild(id, player, name string) (*http.Response, error) {
	return c.Send(http.MethodGet, c.Authentication(), "guild", &Params{
		"id":     id,
		"player": player,
		"name":   name,
	})
}

// GetGamesInformation Game Information
// Returns information about Hypixel Games. This endpoint is in early development and we are working to add more information when possible
//
// https://api.hypixel.net/#tag/Resources/paths/~1v2~1resources~1games/get
func (c *Client) GetGamesInformation() (*http.Response, error) {
	return c.Send(http.MethodGet, nil, "resources/games", nil)
}

// GetAchievements Achievements
//
// https://api.hypixel.net/#tag/Resources/paths/~1v2~1resources~1achievements/get
func (c *Client) GetAchievements() (*http.Response, error) {
	return c.Send(http.MethodGet, nil, "resources/achievements", nil)
}

// GetChallenges Challenges
//
// https://api.hypixel.net/#tag/Resources/paths/~1v2~1resources~1challenges/get
func (c *Client) GetChallenges() (*http.Response, error) {
	return c.Send(http.MethodGet, nil, "resources/challenges", nil)
}

// GetQuests Quests
//
// https://api.hypixel.net/#tag/Resources/paths/~1v2~1resources~1quests/get
func (c *Client) GetQuests() (*http.Response, error) {
	return c.Send(http.MethodGet, nil, "resources/quests", nil)
}

// GetGuildAchievements Guild Achievements
//
// https://api.hypixel.net/#tag/Resources/paths/~1v2~1resources~1guilds~1achievements/get
func (c *Client) GetGuildAchievements() (*http.Response, error) {
	return c.Send(http.MethodGet, nil, "resources/guilds/achievements", nil)
}

// GetVanityPets Vanity Pets
//
// https://api.hypixel.net/#tag/Resources/paths/~1v2~1resources~1vanity~1pets/get
func (c *Client) GetVanityPets() (*http.Response, error) {
	return c.Send(http.MethodGet, nil, "resources/vanity/pets", nil)
}

// GetVanityCompanions Vanity Companions
//
// https://api.hypixel.net/#tag/Resources/paths/~1v2~1resources~1vanity~1companions/get
func (c *Client) GetVanityCompanions() (*http.Response, error) {
	return c.Send(http.MethodGet, nil, "resources/vanity/companions", nil)
}

// GetSkyBlockCollections Collections
// Information regarding Collections in the SkyBlock game.
//
// https://api.hypixel.net/#tag/SkyBlock/paths/~1v2~1resources~1skyblock~1collections/get
func (c *Client) GetSkyBlockCollections() (*http.Response, error) {
	return c.Send(http.MethodGet, nil, "resources/skyblock/collections", nil)
}

// GetSkyBlockSkills Skills
// Information regarding skills in the SkyBlock game.
//
// https://api.hypixel.net/#tag/SkyBlock/paths/~1v2~1resources~1skyblock~1skills/get
func (c *Client) GetSkyBlockSkills() (*http.Response, error) {
	return c.Send(http.MethodGet, nil, "resources/skyblock/skills", nil)
}

// GetSkyBlockItems Items
// Information regarding items in the SkyBlock game.
//
// https://api.hypixel.net/#tag/SkyBlock/paths/~1v2~1resources~1skyblock~1items/get
func (c *Client) GetSkyBlockItems() (*http.Response, error) {
	return c.Send(http.MethodGet, nil, "resources/skyblock/items", nil)
}

// GetSkyBlockElectionAndMayor Election and Mayor
// Information regarding the current mayor and ongoing election in SkyBlock.
//
// https://api.hypixel.net/#tag/SkyBlock/paths/~1v2~1resources~1skyblock~1election/get
func (c *Client) GetSkyBlockElectionAndMayor() (*http.Response, error) {
	return c.Send(http.MethodGet, nil, "resources/skyblock/election", nil)
}

// GetSkyBlockCurrentBingoEvent Current Bingo Event
// Information regarding the current bingo event and its goals.
//
// https://api.hypixel.net/#tag/SkyBlock/paths/~1v2~1resources~1skyblock~1bingo/get
func (c *Client) GetSkyBlockCurrentBingoEvent() (*http.Response, error) {
	return c.Send(http.MethodGet, nil, "resources/skyblock/bingo", nil)
}

// GetSkyBlockNews News
// NEED API Key
//
// https://api.hypixel.net/#tag/SkyBlock/paths/~1v2~1resources~1skyblock~1news/get
func (c *Client) GetSkyBlockNews() (*http.Response, error) {
	return c.Send(http.MethodGet, c.Authentication(), "skyblock/news", nil)
}

// GetAuctions Request auction(s) by the auction UUID, player UUID, or profile UUID.
// Returns the auctions selected by the provided query. Only one query parameter can be used in a single request, and cannot be filtered by multiple.
// NEED API Key
//
// https://api.hypixel.net/#tag/SkyBlock/paths/~1v2~1skyblock~1auction/get
func (c *Client) GetAuctions(uuid, player, profile string) (*http.Response, error) {
	return c.Send(http.MethodGet, c.Authentication(), "skyblock/auction", &Params{
		"uuid":    uuid,
		"player":  player,
		"profile": profile,
	})
}

// GetActiveAuctions Active auctions
// Returns the currently active auctions sorted by last updated first and paginated.
//
// https://api.hypixel.net/#tag/SkyBlock/paths/~1v2~1skyblock~1auctions/get
func (c *Client) GetActiveAuctions(page uint) (*http.Response, error) {
	return c.Send(http.MethodGet, nil, "skyblock/auctions", &Params{
		"page": page,
	})
}

// GetRecentlyEndedAuctions Recently ended auctions
// SkyBlock auctions which ended in the last 60 seconds.
//
// https://api.hypixel.net/#tag/SkyBlock/paths/~1v2~1skyblock~1auctions_ended/get
func (c *Client) GetRecentlyEndedAuctions() (*http.Response, error) {
	return c.Send(http.MethodGet, nil, "skyblock/auctions_ended", nil)
}

// GetBazaar Bazaar
// Returns the list of products along with their sell summary, buy summary and quick status.
// Product Description
// The returned product info has 3 main fields:
//
// buy_summary
// sell_summary
// quick_status
// buy_summary and are the current top 30 orders for each transaction type (in-game example: Stock of Stonks).sell_summary
//
// quick_status is a computed summary of the live state of the product (used for advanced mode view in the bazaar):
//
// sellVolume and are the sum of item amounts in all orders.buyVolume
// sellPrice and are the weighted average of the top 2% of orders by volume.buyPrice
// movingWeek is the historic transacted volume from last 7d + live state.
// sellOrders and are the count of active orders. buyOrders
func (c *Client) GetBazaar() (*http.Response, error) {
	return c.Send(http.MethodGet, nil, "skyblock/bazaar", nil)
}

// GetProfileByUUID Profile by UUID
// SkyBlock profile data, such as stats, objectives etc. The data returned can differ depending on the players in-game API settings.
// NEED API Key
//
// https://api.hypixel.net/#tag/SkyBlock/paths/~1v2~1skyblock~1profile/get
func (c *Client) GetProfileByUUID(profile string) (*http.Response, error) {
	return c.Send(http.MethodGet, c.Authentication(), "skyblock/profile", &Params{
		"profile": profile,
	})
}

// GetProfilesByPlayer Profiles by player
// NEED API Key
//
// https://api.hypixel.net/#tag/SkyBlock/paths/~1v2~1skyblock~1profiles/get
func (c *Client) GetProfilesByPlayer(uuid string) (*http.Response, error) {
	return c.Send(http.MethodGet, c.Authentication(), "skyblock/profiles", &Params{
		"uuid": uuid,
	})
}

// GetMuseumData Museum data by profile ID
// SkyBlock museum data for all members of the provided profile. The data returned can differ depending on the players in-game API settings.
// NEED API Key
//
// https://api.hypixel.net/#tag/SkyBlock/paths/~1v2~1skyblock~1museum/get
func (c *Client) GetMuseumData(profile string) (*http.Response, error) {
	return c.Send(http.MethodGet, c.Authentication(), "skyblock/museum", &Params{
		"profile": profile,
	})
}

// GetGardenData Garden data by profile ID
// SkyBlock garden data for the provided profile.
// NEED API Key
//
// https://api.hypixel.net/#tag/SkyBlock/paths/~1v2~1skyblock~1garden/get
func (c *Client) GetGardenData(profile string) (*http.Response, error) {
	return c.Send(http.MethodGet, c.Authentication(), "skyblock/garden", &Params{
		"profile": profile,
	})
}

// GetBingoData Bingo data by player
// Bingo data for participated events of the provided player.
// NEED API Key
//
// https://api.hypixel.net/#tag/SkyBlock/paths/~1v2~1skyblock~1bingo/get
func (c *Client) GetBingoData(uuid string) (*http.Response, error) {
	return c.Send(http.MethodGet, c.Authentication(), "bingo", &Params{
		"uuid": uuid,
	})
}

// GetActiveOrUpcomingFireSales Active/Upcoming Fire Sales
// Retrieve the currently active or upcoming Fire Sales for SkyBlock.
//
// https://api.hypixel.net/#tag/SkyBlock/paths/~1v2~1skyblock~1firesales/get
func (c *Client) GetActiveOrUpcomingFireSales() (*http.Response, error) {
	return c.Send(http.MethodGet, nil, "skyblock/firesales", nil)
}

// GetCurrentlyActivePublicHouses currently active public houses.
// This data may be cached for a short period of time.
// NEED API Key
//
// https://api.hypixel.net/#tag/Housing/paths/~1v2~1housing~1active/get
func (c *Client) GetCurrentlyActivePublicHouses() (*http.Response, error) {
	return c.Send(http.MethodGet, c.Authentication(), "housing/active", nil)
}

// GetSpecificHouseInformation Information about a specific house.
// This data may be cached for a short period of time.
// NEED API Key
//
// https://api.hypixel.net/#tag/Housing/paths/~1v2~1housing~1house/get
func (c *Client) GetSpecificHouseInformation(house string) (*http.Response, error) {
	return c.Send(http.MethodGet, c.Authentication(), "housing/house", &Params{
		"house": house,
	})
}

// GetSpecificPlayerPublicHouses The public houses for a specific player.//
// This data may be cached for a short period of time.
// NEED API Key
//
// https://api.hypixel.net/#tag/Housing/paths/~1v2~1housing~1houses/get
func (c *Client) GetSpecificPlayerPublicHouses(player string) (*http.Response, error) {
	return c.Send(http.MethodGet, c.Authentication(), "housing/houses", &Params{
		"player": player,
	})
}

// GetActiveNetworkBoosters Active Network Boosters
// NEED API Key
//
// https://api.hypixel.net/#tag/Other/paths/~1v2~1boosters/get
func (c *Client) GetActiveNetworkBoosters() (*http.Response, error) {
	return c.Send(http.MethodGet, c.Authentication(), "boosters", nil)
}

// GetCurrentPlayerCounts Current Player Counts
// NEED API Key
//
// https://api.hypixel.net/#tag/Other/paths/~1v2~1counts/get
func (c *Client) GetCurrentPlayerCounts() (*http.Response, error) {
	return c.Send(http.MethodGet, c.Authentication(), "counts", nil)
}

// GetCurrentLeaderboards Current Leaderboards
// NEED API Key
//
// https://api.hypixel.net/#tag/Other/paths/~1v2~1leaderboards/get
func (c *Client) GetCurrentLeaderboards() (*http.Response, error) {
	return c.Send(http.MethodGet, c.Authentication(), "leaderboards", nil)
}

// GetPunishmentStatistics Punishment Statistics
// NEED API Key
//
// https://api.hypixel.net/#tag/Other/paths/~1v2~1punishmentstats/get
func (c *Client) GetPunishmentStatistics() (*http.Response, error) {
	return c.Send(http.MethodGet, c.Authentication(), "punishmentstats", nil)
}
