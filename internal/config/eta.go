// Feature config
package config

type EtaConfig struct {
	// Max number of entities allowed in fetch Eta request body
	MaxEntities int `mapstructure:"maxEntities"`
	// Buffer added to the legacy OMS-style ETA range (seconds).
	// When unset/<=0, constants.DefaultLegacyEtaBufferSeconds is used.
	LegacyEtaBufferInSeconds int `mapstructure:"legacyEtaBufferInSeconds"`
}
