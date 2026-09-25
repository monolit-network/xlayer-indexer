package models

type NotifierEvents string

const (
	NotifierEventCondPrepared   NotifierEvents = "condition_prepared_topic"
	NotifierEventCondResolved   NotifierEvents = "condition_resolved_topic"
	NotifierEventAdapterUMA     NotifierEvents = "adapter_uma_topic"
	NotifierEventOrder          NotifierEvents = "order_topic"
	NotifierEventReorg          NotifierEvents = "reorg_topic"
	NotifierEventSolanaSwap     NotifierEvents = "solana_swap_topic"
	NotifierEventSolanaTransfer NotifierEvents = "solana_transfer_topic"
	NotifierEventSolanaMint     NotifierEvents = "solana_mint_topic"
	NotifierEventSolanaBurn     NotifierEvents = "solana_burn_topic"
	NotifierEventEVMSwap        NotifierEvents = "evm_swap_topic"
	NotifierEventEVMTransfer    NotifierEvents = "evm_transfer_topic"
	NotifierEventEVMDefi        NotifierEvents = "evm_defi_topic"
)
