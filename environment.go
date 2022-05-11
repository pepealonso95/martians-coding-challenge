package main

import (
	"strconv"
)

// Representation of a robot
type Environment struct {

	// Lost robots are stored in a map to make it as fast as possible to check if a robot has been lost
	// A matrix is possible too by knowing the max limits of 50x50, however its not as efficient in memory
	// and potentially applied to a real world scenario, storing every single position in mars
	// just to check lost robots would be far too much to handle

	// Map storing all lost robot positions
	lostRobots map[string]bool
	// X coordinate of environment limit
	limitX int
	// Y coordinate of environment limit
	limitY int
}

// Check if the coordinates have a lost robot
func (env *Environment) checkLostRobot(x int, y int) bool {
	pos := strconv.Itoa(x) + " " + strconv.Itoa(y)
	if val, ok := env.lostRobots[pos]; ok {
		return val
	}
	return false
}
