package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
)

// Storing limits and lost robots as global variables is not the best way to do this
// when the program is running in a container or scales, but it works for now

var limitX int
var limitY int

// Lost robots are stored in a map to make it as fast as possible to check if a robot has been lost
// A matrix is possible too by knowing the max limits of 50x50, however its not as efficient in memory
// and potentially applied to a real world scenario, storing every single position in mars
// just to check lost robots would be far too much to handle

var lostRobots map[string]bool

// The struct and methods for robots should probably be in a separate class,
// but due to the narrow scope of the project, it was kept here

// Representation of a robot
type Robot struct {
	// Whether robot is lost or not
	Lost bool
	// Where the robot is pointing to (N,E,S,W)
	Orientation string
	// Latitude coordinate representation of robot position
	X int
	// Longitude coordinate representation of robot position
	Y int
}

// Returns a string representation of the robot's position
func (robot *Robot) GetPos() string {
	return strconv.Itoa(robot.X) + " " + strconv.Itoa(robot.Y)
}

// Decides how to process the given instruction for the robot
func (robot *Robot) Instruct(direction string) (bool, error) {
	switch direction {
	case "F":
		return robot.Move()
	case "L":
		return robot.TurnLeft()
	case "R":
		return robot.TurnRight()
	default:
		return false, errors.New("Invalid instruction")
	}
}

// NOTE: There are smarter ways to do turn left and right, but the simplest possible (and most readable)
// are these two switch statement methods (two separate methods as code is not really being that reused here)

// In GO there is no generic library function to get index of value in slice, so
// getting the position with the indes +1 or -1 depending on the direction is not as simple as in say JS
// In JS it could be done using an ordered array of orientations [N,E,S,W] and doing:
// orientations[(orientations.indexOf(robot.orientation) + 1 (RIGHT) or -1 (LEFT) + 4) % 4]

// It doesnt matter that much though as the orientations are only 4 so the code looks clean enough

// Function to turn right
func (robot *Robot) TurnRight() (bool, error) {
	switch robot.Orientation {
	case "N":
		robot.Orientation = "E"
		return false, nil
	case "E":
		robot.Orientation = "S"
		return false, nil
	case "W":
		robot.Orientation = "N"
		return false, nil
	case "S":
		robot.Orientation = "W"
		return false, nil
	default:
		return true, errors.New("Invalid orientation")
	}
}

// Function to turn left
func (robot *Robot) TurnLeft() (bool, error) {
	switch robot.Orientation {
	case "N":
		robot.Orientation = "W"
		return false, nil
	case "E":
		robot.Orientation = "N"
		return false, nil
	case "W":
		robot.Orientation = "S"
		return false, nil
	case "S":
		robot.Orientation = "E"
		return false, nil
	default:
		return true, errors.New("Invalid orientation")
	}
}

// Checks in which direction and lets the robot advance if no lost robot is on its way
func (robot *Robot) Move() (bool, error) {
	switch robot.Orientation {
	case "N":
		if checkLostRobot(robot.X, robot.Y+1) {
			return false, nil
		}
		return robot.advance(robot.X, robot.Y+1), nil
	case "E":
		if checkLostRobot(robot.X+1, robot.Y) {
			return false, nil
		}
		return robot.advance(robot.X+1, robot.Y), nil
	case "W":
		if checkLostRobot(robot.X-1, robot.Y) {
			return false, nil
		}
		return robot.advance(robot.X-1, robot.Y), nil
	case "S":
		if checkLostRobot(robot.X, robot.Y-1) {
			return false, nil
		}
		return robot.advance(robot.X, robot.Y-1), nil
	default:
		return true, errors.New("Invalid orientation")
	}
}

func (robot *Robot) printRobotPos() string {
	robotLostString := ""
	if robot.Lost {
		robotLostString = "LOST"
	}
	return strconv.Itoa(robot.X) + " " + strconv.Itoa(robot.Y) + " " + robot.Orientation + " " + robotLostString
}

// Advance the robot to the next position or report it as lost and keep its last seen position
func (robot *Robot) advance(x int, y int) bool {
	if y < 0 || x < 0 || y > limitY || x > limitX {
		robot.Lost = true
		lostRobots[strconv.Itoa(x)+" "+strconv.Itoa(y)] = true
		return true
	} else {
		robot.X = x
		robot.Y = y
		return false
	}
}

// Check if the coordinates have a lost robot
func checkLostRobot(x int, y int) bool {
	pos := strconv.Itoa(x) + " " + strconv.Itoa(y)
	if val, ok := lostRobots[pos]; ok {
		return val
	}
	return false
}

// Scan and set limits from the input text file
func setInputLimits(scanner *bufio.Scanner) {
	scanner.Scan()

	limits := scanner.Text()

	x, _ := strconv.Atoi(limits[0:1])
	y, _ := strconv.Atoi(limits[2:3])

	if x < 0 || x > 50 || y < 0 || y > 50 {
		log.Fatal(errors.New("Invalid limits"))
	}

	limitX = x
	limitY = y
}

// Get a robot from the input text position string
func getInputRobot(robotPos string) Robot {
	x, _ := strconv.Atoi(robotPos[0:1])
	y, _ := strconv.Atoi(robotPos[2:3])
	orientation := robotPos[4:5]

	return Robot{X: x, Y: y, Orientation: orientation, Lost: false}
}

// Takes the input string instructions, makes the robot process them and prints the result
func processInputInstruction(instructions string, robot *Robot) {
	chars := []rune(instructions)
	gotLost := false
	for i := 0; i < len(chars) && !robot.Lost && !gotLost; i++ {
		char := string(chars[i])
		var err error
		gotLost, err = robot.Instruct(char)
		if err != nil {
			log.Fatal(err)
			gotLost = true
		}
	}
	fmt.Println(robot.printRobotPos())

	// Non-Lost Robots are not stored anywhere as colisions are not contemplated by the problem
	// so persisting their position is not necessary
	// If colisions were contemplated, storing them just like lost robots would be the way to go
}

func main() {

	var filePath string

	fmt.Println("Enter the file path with its format: (pressing enter defaults to 'instructions.txt')")
	fmt.Scanf("%s", &filePath)

	if filePath == "" {
		filePath = "instructions.txt"
	}

	file, err := os.Open(filePath)

	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	setInputLimits(scanner)

	lostRobots = make(map[string]bool)

	// Loop while there is a line to read
	for scanner.Scan() {
		robotPos := scanner.Text()
		robot := getInputRobot(robotPos)

		scanner.Scan()
		instructions := scanner.Text()
		// Speed of the algorithm might be improved by running the instructions in gourutines,
		// however the expected output suggests that the order of instructions are important
		processInputInstruction(instructions, &robot)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

}
