package queries

const (
	CreateEvmParserRegistryTableSQL = `
		CREATE TABLE IF NOT EXISTS evm_parser_registry (
			chain VARCHAR(255) NOT NULL,
			contract_address VARCHAR(255) NOT NULL,
			instruction_hash VARCHAR(255) NOT NULL,
			parser_name VARCHAR(255) NOT NULL,
			source VARCHAR(255) NOT NULL,
			PRIMARY KEY (chain, contract_address, instruction_hash)
		);

		CREATE INDEX IF NOT EXISTS evm_parser_registry_chain_idx ON evm_parser_registry (chain);

		CREATE INDEX IF NOT EXISTS evm_parser_registry_contract_address_idx ON evm_parser_registry (contract_address);

		CREATE INDEX IF NOT EXISTS evm_parser_registry_instruction_hash_idx ON evm_parser_registry (instruction_hash);

		CREATE INDEX IF NOT EXISTS evm_parser_registry_parser_name_idx ON evm_parser_registry (parser_name);`

	UpsertEvmParserSQL = `
		INSERT INTO evm_parser_registry (chain, contract_address, instruction_hash, parser_name, source)
		VALUES (LOWER($1), LOWER($2), LOWER($3), LOWER($4), LOWER($5))
		ON CONFLICT (chain, contract_address, instruction_hash)
		DO UPDATE SET parser_name = LOWER(EXCLUDED.parser_name), source = LOWER(EXCLUDED.source);`

	GetEvmParserSQL = `
		SELECT parser_name, source FROM evm_parser_registry WHERE LOWER(chain) = LOWER($1) AND LOWER(contract_address) = LOWER($2) AND LOWER(instruction_hash) = LOWER($3);`

	GetAllEvmParsersForChainSQL = `
		SELECT LOWER(contract_address), LOWER(instruction_hash), LOWER(parser_name), LOWER(source) FROM evm_parser_registry WHERE LOWER(chain) = LOWER($1);`
)
