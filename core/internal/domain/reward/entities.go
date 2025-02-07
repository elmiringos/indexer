package reward

type Reward struct {
	BlockHash   string
	AddressHash string
	Amount      string
}

func (r *Reward) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"block_hash":   r.BlockHash,
		"address_hash": r.AddressHash,
		"amount":       r.Amount,
	}
}

func MakeRewardSlice(rewards []*Reward) []map[string]interface{} {
	slices := make([]map[string]interface{}, len(rewards))
	for i, reward := range rewards {
		slices[i] = reward.ToMap()
	}
	return slices
}
