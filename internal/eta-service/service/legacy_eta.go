package service

import "github.com/nutanalabs/eta-service/internal/constants"

// calculateLegacyEta replicates order-management-service's cart ETA formula
// (calculateDeliveryEstimates + buildCartMeta) using data eta-service
// already has, purely for shadow comparison. It is intentionally isolated
// from the primary FetchEta formula and must never be called from it.
func calculateLegacyEta(rat, kpt, cat, fm, pickup, lastMile float64, bufferSeconds int) (etaSeconds uint, displayMin uint, displayMax uint) {
	total := rat + max(kpt, cat+fm+pickup) + lastMile

	buffer := bufferSeconds
	if buffer <= 0 {
		buffer = constants.DefaultLegacyEtaBufferSeconds
	}

	minVal := uint(total)
	maxVal := uint(total) + uint(buffer)
	return minVal, minVal, maxVal
}
