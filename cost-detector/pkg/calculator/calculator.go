package calculator

import "cost-detector/pkg/models"

// Calculator does the cost calculations
type Calculator struct {
	NodePrices map[string]float64 // Map of instance type to cost per hour
}

// NewCalculator creates a new cost calculator
func NewCalculator() *Calculator {
	return &Calculator{
		NodePrices: make(map[string]float64),
	}
}

// CalculatePodCost calculates hourly cost of a pod
func (c *Calculator) CalculatePodCost(pod *models.Pod) float64 {
	// This is simplified. In reality:
	// Cost = (Pod CPU Request / Node CPU) * Node Price Per Hour
	//      + (Pod Memory Request / Node Memory) * Node Price Per Hour
	
	// For now, rough estimate: $0.05 per CPU per hour, $0.01 per GB per hour
	cpuCost := pod.CPU * 0.05
	memoryCost := pod.Memory * 0.01
	
	return cpuCost + memoryCost
}

// CalculateHourlyCost calculates cost for multiple pods
func (c *Calculator) CalculateHourlyCost(pods []*models.Pod) float64 {
	totalCost := 0.0
	for _, pod := range pods {
		totalCost += c.CalculatePodCost(pod)
	}
	return totalCost
}
