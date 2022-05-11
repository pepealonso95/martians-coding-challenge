package main

import (
	"strconv"
)

// Representation of the environment
type Environment struct {

	// Lost robots are stored in a map to make it as fast as possible to check if a robot has been lost
	// A matrix is possible too by knowing the max limits of 50x50, however its not as efficient in memory
	// and potentially applied to a real world scenario, storing every single position in mars
	// just to check lost robots would be far too much to handle
	// The map uses the type struct{} to store whether there is a robot or not instead of bool
	// because using an empty struct occupies the least amount of memory possible (compared to bool, string, etc)

	// Map storing all lost robot positions
	lostRobots map[string]struct{}
	// X coordinate of environment limit
	limitX int
	// Y coordinate of environment limit
	limitY int
}

// Check if the coordinates have a lost robot
func (env *Environment) checkLostRobot(x int, y int) bool {
	pos := strconv.Itoa(x) + " " + strconv.Itoa(y)
	if _, ok := env.lostRobots[pos]; ok {
		return true
	}
	return false
}

// Check if the coordinates have a lost robot
func (env *Environment) addLostRobot(x int, y int) {
	env.lostRobots[strconv.Itoa(x)+" "+strconv.Itoa(y)] = struct{}{}
}
