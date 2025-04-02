package smartcontract

type SmartContract struct {
	AddressHash      string
	Name             string
	CompilerVersion  string
	SourceCode       string
	ABI              string
	CompilerSettings string
	VerifiedByEth    bool
	EvmVersion       string
}

func (s *SmartContract) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"address_hash":      s.AddressHash,
		"name":              s.Name,
		"compiler_version":  s.CompilerVersion,
		"source_code":       s.SourceCode,
		"abi":               s.ABI,
		"compiler_settings": s.CompilerSettings,
		"verified_by_eth":   s.VerifiedByEth,
		"evm_version":       s.EvmVersion,
	}
}

func MakeSlice(smartContracts []*SmartContract) []map[string]interface{} {
	slices := make([]map[string]interface{}, len(smartContracts))
	for i, smartContract := range smartContracts {
		slices[i] = smartContract.ToMap()
	}
	return slices
}
