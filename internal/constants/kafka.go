package constants

const (
	FETCH_ETA_EVENTS_TOPIC string = "domain.geo.eta_estimates_events"
)

// DefaultLegacyEtaBufferSeconds mirrors order-management-service's
// DefaultBufferTimeForETA (5 minutes), converted to seconds, used when
// EtaConfig.LegacyEtaBufferInSeconds is unset.
const DefaultLegacyEtaBufferSeconds = 300
