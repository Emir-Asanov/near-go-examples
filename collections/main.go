package main

import "github.com/vlmoon99/near-sdk-go/collections"

// ── Instantiation ─────────────────────────────────────────────────────────────

// @contract:state
type CollectionsContract struct {
	MyVector          *collections.Vector[string]               `json:"my_vector"`
	MyLookupMap       *collections.LookupMap[string, string]    `json:"my_lookup_map"`
	MyUnorderedMap    *collections.UnorderedMap[string, string] `json:"my_unordered_map"`
	MyLookupSet       *collections.LookupSet[string]            `json:"my_lookup_set"`
	MyUnorderedSet    *collections.UnorderedSet[string]         `json:"my_unordered_set"`
	MyTreeMap         *collections.TreeMap[string, uint64]      `json:"my_tree_map"`
	UserAssetPrefixes *collections.LookupMap[string, string]    `json:"user_asset_prefixes"`
	UserAssetLengths  *collections.LookupMap[string, uint64]    `json:"user_asset_lengths"`
}

// @contract:init
func (c *CollectionsContract) Init() {
	c.MyVector          = collections.NewVector[string]("v")
	c.MyLookupMap       = collections.NewLookupMap[string, string]("m")
	c.MyUnorderedMap    = collections.NewUnorderedMap[string, string]("um")
	c.MyLookupSet       = collections.NewLookupSet[string]("s")
	c.MyUnorderedSet    = collections.NewUnorderedSet[string]("us")
	c.MyTreeMap         = collections.NewTreeMap[string, uint64]("tm")
	c.UserAssetPrefixes = collections.NewLookupMap[string, string]("up")
	c.UserAssetLengths  = collections.NewLookupMap[string, uint64]("ul")
}

// ── Vector ────────────────────────────────────────────────────────────────────

// @contract:mutating
func (c *CollectionsContract) VectorPush(value string) uint64 {
	c.MyVector.Push(value)
	return c.MyVector.Length()
}

// @contract:view
func (c *CollectionsContract) VectorGet(index uint64) string {
	val, err := c.MyVector.Get(index)
	if err != nil {
		return ""
	}
	return val
}

// @contract:view
func (c *CollectionsContract) VectorGetAll() []string {
	result, _ := c.MyVector.ToSlice()
	return result
}

// @contract:mutating
func (c *CollectionsContract) VectorPop() string {
	val, _ := c.MyVector.Pop()
	return val
}

// ── LookupMap ─────────────────────────────────────────────────────────────────

// @contract:mutating
func (c *CollectionsContract) LookupMapSet(key string, value string) {
	c.MyLookupMap.Insert(key, value)
}

// @contract:view
func (c *CollectionsContract) LookupMapGet(key string) string {
	val, err := c.MyLookupMap.Get(key)
	if err != nil {
		return ""
	}
	return val
}

// @contract:view
func (c *CollectionsContract) LookupMapContains(key string) bool {
	ok, _ := c.MyLookupMap.Contains(key)
	return ok
}

// @contract:mutating
func (c *CollectionsContract) LookupMapRemove(key string) {
	c.MyLookupMap.Remove(key)
}

// ── UnorderedMap ──────────────────────────────────────────────────────────────

// @contract:mutating
func (c *CollectionsContract) UnorderedMapSet(key string, value string) {
	c.MyUnorderedMap.Insert(key, value)
}

// @contract:view
func (c *CollectionsContract) UnorderedMapGet(key string) string {
	val, err := c.MyUnorderedMap.Get(key)
	if err != nil {
		return ""
	}
	return val
}

// @contract:view
func (c *CollectionsContract) UnorderedMapKeys() []string {
	keys, _ := c.MyUnorderedMap.Keys()
	return keys
}

// @contract:view
func (c *CollectionsContract) UnorderedMapValues() []string {
	values, _ := c.MyUnorderedMap.Values()
	return values
}

// @contract:mutating
func (c *CollectionsContract) UnorderedMapRemove(key string) {
	c.MyUnorderedMap.Remove(key)
}

// ── LookupSet ─────────────────────────────────────────────────────────────────

// @contract:mutating
func (c *CollectionsContract) LookupSetAdd(value string) {
	c.MyLookupSet.Insert(value)
}

// @contract:view
func (c *CollectionsContract) LookupSetContains(value string) bool {
	ok, _ := c.MyLookupSet.Contains(value)
	return ok
}

// @contract:mutating
func (c *CollectionsContract) LookupSetRemove(value string) {
	c.MyLookupSet.Remove(value)
}

// ── UnorderedSet ──────────────────────────────────────────────────────────────

// @contract:mutating
func (c *CollectionsContract) UnorderedSetAdd(value string) {
	c.MyUnorderedSet.Insert(value)
}

// @contract:view
func (c *CollectionsContract) UnorderedSetContains(value string) bool {
	ok, _ := c.MyUnorderedSet.Contains(value)
	return ok
}

// @contract:view
func (c *CollectionsContract) UnorderedSetAll() []string {
	all, _ := c.MyUnorderedSet.All()
	return all
}

// @contract:mutating
func (c *CollectionsContract) UnorderedSetRemove(value string) {
	c.MyUnorderedSet.Remove(value)
}

// ── TreeMap ───────────────────────────────────────────────────────────────────

// @contract:mutating
func (c *CollectionsContract) TreeMapSet(key string, score uint64) {
	c.MyTreeMap.Insert(key, score)
}

// @contract:view
func (c *CollectionsContract) TreeMapGetMin() string {
	key, _ := c.MyTreeMap.MinKey()
	return key
}

// @contract:view
func (c *CollectionsContract) TreeMapGetMax() string {
	key, _ := c.MyTreeMap.MaxKey()
	return key
}

// @contract:view
func (c *CollectionsContract) TreeMapGetAllKeys() []string {
	keys, _ := c.MyTreeMap.Keys()
	return keys
}

// ── Pagination ────────────────────────────────────────────────────────────────

// @contract:view
func (c *CollectionsContract) GetPage(fromIndex uint64, limit uint64) []string {
	total := c.MyVector.Length()
	if fromIndex >= total {
		return []string{}
	}
	end := fromIndex + limit
	if end > total {
		end = total
	}
	result := make([]string, 0, end-fromIndex)
	for i := fromIndex; i < end; i++ {
		val, err := c.MyVector.Get(i)
		if err == nil {
			result = append(result, val)
		}
	}
	return result
}

// ── Nested Collections ────────────────────────────────────────────────────────

// @contract:mutating
func (c *CollectionsContract) AddAsset(userId string, asset string) {
	prefix := "assets:" + userId
	c.UserAssetPrefixes.Insert(userId, prefix)

	currentLen, _ := c.UserAssetLengths.Get(userId)
	userAssets := &collections.Vector[string]{Prefix: prefix, Len: currentLen}
	userAssets.Push(asset)

	c.UserAssetLengths.Insert(userId, userAssets.Length())
}

// @contract:view
func (c *CollectionsContract) GetUserAssets(userId string) []string {
	prefix, err := c.UserAssetPrefixes.Get(userId)
	if err != nil {
		return []string{}
	}
	currentLen, _ := c.UserAssetLengths.Get(userId)
	userAssets := &collections.Vector[string]{Prefix: prefix, Len: currentLen}
	result, _ := userAssets.ToSlice()
	return result
}
