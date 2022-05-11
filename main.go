package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
)

// Position of x coordinate in input string
const xCoordPos = 0

// Position of y coordinate in input string
const yCoordPos = 2

// Position of orientation in input string
const orientationPos = 4

// Max amount of cells in the map
const mapLimit = 50

// Scan and set limits from the input text file
func createEnvWithLimits(x int, y int) Environment {

	if x < 0 || x > mapLimit || y < 0 || y > mapLimit {
		log.Fatal(errors.New("Invalid limits"))
	}

	env := Environment{lostRobots: make(map[string]struct{}), limitX: x, limitY: y}

	return env

}

// Get a robot from the input text position string
func getInputRobot(robotPos string) Robot {
	x, _ := strconv.Atoi(string(robotPos[xCoordPos]))
	y, _ := strconv.Atoi(string(robotPos[yCoordPos]))
	orientation := string(robotPos[orientationPos])
	robot := Robot{X: x, Y: y, Lost: false}
	robot.setOrientation(orientation)
	return robot
}

// Takes the input string instructions, makes the robot process them and prints the result
func processInputInstruction(instructions string, robot *Robot, env *Environment) {
	gotLost := false
	for i := 0; i < len(instructions) && !robot.Lost && !gotLost; i++ {
		char := string(instructions[i])
		var err error
		gotLost, err = robot.Instruct(char, env)
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

	scanner.Scan()

	limits := scanner.Text()

	x, _ := strconv.Atoi(string(limits[xCoordPos]))
	y, _ := strconv.Atoi(string(limits[yCoordPos]))

	env := createEnvWithLimits(x, y)

	// Loop while there is a line to read
	for scanner.Scan() {
		robotPos := scanner.Text()
		robot := getInputRobot(robotPos)

		scanner.Scan()
		instructions := scanner.Text()
		// Speed of the algorithm might be improved by running the instructions in gourutines,
		// however the expected output suggests that the order of instructions are important
		processInputInstruction(instructions, &robot, &env)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

}
