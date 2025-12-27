package tacticus_wrapper

type TacticusPlayerResponse struct {
	Player   PlayerData `json:"player"`
	MetaData MetaData   `json:"metaData"`
}

type PlayerData struct {
	Details   PlayerDetails   `json:"details"`
	Units     []PlayerUnit    `json:"units"`
	Inventory PlayerInventory `json:"inventory"`
	Progress  PlayerProgress  `json:"progress"`
}

type PlayerDetails struct {
	Name       string `json:"name"`
	PowerLevel int    `json:"powerLevel"`
}

type PlayerUnit struct {
	ID               string     `json:"id"`
	Name             string     `json:"name"`
	Faction          string     `json:"faction"`
	GrandAlliance    string     `json:"grandAlliance"`
	ProgressionIndex int        `json:"progressionIndex"`
	XP               int        `json:"xp"`
	XPLevel          int        `json:"xpLevel"`
	Rank             int        `json:"rank"`
	Abilities        []Ability  `json:"abilities"`
	Upgrades         []int      `json:"upgrades"`
	Items            []UnitItem `json:"items"`
	Shards           int        `json:"shards"`
	MythicShards     int        `json:"mythicShards"`
}

type Ability struct {
	ID    string `json:"id"`
	Level int    `json:"level"`
}

type UnitItem struct {
	SlotID string `json:"slotId"`
	Level  int    `json:"level"`
	ID     string `json:"id"`
	Name   string `json:"name"`
	Rarity string `json:"rarity"`
}

type PlayerInventory struct {
	Items             []InventoryItem           `json:"items"`
	Upgrades          []InventoryUpgrade        `json:"upgrades"`
	Shards            []ShardInfo               `json:"shards"`
	MythicShards      []ShardInfo               `json:"mythicShards"`
	XPBooks           []XPBook                  `json:"xpBooks"`
	AbilityBadges     map[string][]AbilityBadge `json:"abilityBadges"`
	Components        []Component               `json:"components"`
	ForgeBadges       []ForgeBadge              `json:"forgeBadges"`
	Orbs              map[string][]Orb          `json:"orbs"`
	RequisitionOrders RequisitionOrders         `json:"requisitionOrders"`
	ResetStones       int                       `json:"resetStones"`
}

type InventoryItem struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Level  int    `json:"level"`
	Amount int    `json:"amount"`
}

type InventoryUpgrade struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Amount int    `json:"amount"`
}

type ShardInfo struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Amount int    `json:"amount"`
}

type XPBook struct {
	ID     string `json:"id"`
	Rarity string `json:"rarity"`
	Amount int    `json:"amount"`
}

type AbilityBadge struct {
	Name   string `json:"name"`
	Rarity string `json:"rarity"`
	Amount int    `json:"amount"`
}

type Component struct {
	Name          string `json:"name"`
	GrandAlliance string `json:"grandAlliance"`
	Amount        int    `json:"amount"`
}

type ForgeBadge struct {
	Name   string `json:"name"`
	Rarity string `json:"rarity"`
	Amount int    `json:"amount"`
}

type Orb struct {
	Rarity string `json:"rarity"`
	Amount int    `json:"amount"`
}

type RequisitionOrders struct {
	Regular int `json:"regular"`
	Blessed int `json:"blessed"`
}

type PlayerProgress struct {
	Campaigns  []Campaign `json:"campaigns"`
	Arena      Arena      `json:"arena"`
	GuildRaid  GuildRaid  `json:"guildRaid"`
	Onslaught  TokensInfo `json:"onslaught"`
	SalvageRun TokensInfo `json:"salvageRun"`
}

type Campaign struct {
	ID      string   `json:"id"`
	Name    string   `json:"name"`
	Type    string   `json:"type"`
	Battles []Battle `json:"battles"`
}

type Battle struct {
	BattleIndex  int `json:"battleIndex"`
	AttemptsLeft int `json:"attemptsLeft"`
	AttemptsUsed int `json:"attemptsUsed"`
}

type Arena struct {
	Tokens TokensInfo `json:"tokens"`
}

type GuildRaid struct {
	Tokens     TokensInfo `json:"tokens"`
	BombTokens TokensInfo `json:"bombTokens"`
}

type TokensInfo struct {
	Current             int `json:"current"`
	Max                 int `json:"max"`
	NextTokenInSeconds  int `json:"nextTokenInSeconds"`
	RegenDelayInSeconds int `json:"regenDelayInSeconds"`
}

type MetaData struct {
	ConfigHash      string   `json:"configHash"`
	ApiKeyExpiresOn int64    `json:"apiKeyExpiresOn"`
	LastUpdatedOn   int64    `json:"lastUpdatedOn"`
	Scopes          []string `json:"scopes"` // <-- было string
}
