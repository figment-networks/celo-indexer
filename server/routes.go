package server

// setupRoutes sets up routes for gin application
func (s *Server) setupRoutes() {
	s.engine.GET("/health", s.handlers.Health.Handle)
	s.engine.GET("/status", s.handlers.GetStatus.Handle)
	s.engine.GET("/block", s.handlers.GetBlockByHeight.Handle)
	s.engine.GET("/block_times/:limit", s.handlers.GetBlockTimes.Handle)
	s.engine.GET("/blocks_summary", s.handlers.GetBlockSummary.Handle)
	s.engine.GET("/transactions", s.handlers.GetTransactionsByHeight.Handle)
	s.engine.GET("/account_details/:address", s.handlers.GetAccountDetails.Handle)
	s.engine.GET("/account/:address", s.handlers.GetAccountByHeight.Handle)
	s.engine.GET("/validator/:address", s.handlers.GetValidatorByAddress.Handle)
	s.engine.GET("/validators/for_min_height/:height", s.handlers.GetValidatorsForMinHeight.Handle)
	s.engine.GET("/validators", s.handlers.GetValidatorsByHeight.Handle)
	s.engine.GET("/validators_summary", s.handlers.GetValidatorSummary.Handle)
	s.engine.GET("/validator_group/:address", s.handlers.GetValidatorGroupByAddress.Handle)
	s.engine.GET("/validator_groups", s.handlers.GetValidatorGroupsByHeight.Handle)
	s.engine.GET("/validator_groups_summary", s.handlers.GetValidatorGroupSummary.Handle)
	s.engine.GET("/system_events/:address", s.handlers.GetSystemEventsForAddress.Handle)
	s.engine.GET("/system_events", s.handlers.GetSystemEvents.Handle)
	s.engine.GET("/proposals", s.handlers.GetProposals.Handle)
	s.engine.GET("/proposals/:proposal_id/activity", s.handlers.GetProposalActivity.Handle)
}
