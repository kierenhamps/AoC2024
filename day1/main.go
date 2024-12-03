package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"regexp"
	"sort"
	"strconv"
)

var (
	ErrNoNumbersLeftInList   = errors.New("no numbers left in list")
	ErrInputCannotBeZero     = errors.New("input cannot be zero")
	ErrInputCannotBeNegative = errors.New("input cannot be negative")
	ErrInvalidInputFormat    = errors.New("invalid input format")
)

type location int

func NewLocation(i int) (location, error) {
	if i == 0 {
		return 0, ErrInputCannotBeZero
	}
	if i < 0 {
		return 0, ErrInputCannotBeNegative
	}
	return location(i), nil
}

func (l location) Distance(r location) int {
	if l > r {
		return int(l - r)
	}
	return int(r - l)
}

func (l location) Int() int {
	return int(l)
}

type LocationList struct {
	list []location
}

func NewLocationList() *LocationList {
	return &LocationList{}
}

func (ll *LocationList) AddLocation(l location) {
	ll.list = append(ll.list, l)
}

func (ll *LocationList) CountMatches(l location) int {
	var count int
	for _, loc := range ll.list {
		if loc == l {
			count++
		}
	}
	return count
}

func (ll *LocationList) Next() location {
	if len(ll.list) == 0 {
		return 0
	}

	// first sort the list
	sort.Slice(ll.list, func(i, j int) bool {
		return ll.list[i] < ll.list[j]
	})
	next := ll.list[0]

	// remove the first element
	ll.list = ll.list[1:]

	return next
}

func (ll *LocationList) Size() int {
	return len(ll.list)
}

func main() {
	// read the inputs
	inputFile, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer inputFile.Close()

	// Create the lists
	leftList, rightList, err := createLists(inputFile)
	if err != nil {
		log.Fatal(err)
	}

	// Create the pairs and calculate the sum of the distances
	sum := sumDistances(leftList, rightList)
	log.Printf("(PART 1) Sum of distances: %d\n", sum)

	// Reset the input file and regenerate lists for part 2
	inputFile.Seek(0, 0)
	leftList, rightList, err = createLists(inputFile)
	if err != nil {
		log.Fatal(err)
	}

	// Calculate sum of similarity score
	similaritySum := sumSimilarities(leftList, rightList)
	log.Printf("(PART 2) Sum of similarity score: %d\n", similaritySum)
}

func createLists(inputFile *os.File) (*LocationList, *LocationList, error) {
	leftList := NewLocationList()
	rightList := NewLocationList()

	log.Println("Reading input file")
	scanner := bufio.NewScanner(inputFile)
	for scanner.Scan() {
		re := regexp.MustCompile(`(\d+)\s+(\d+)`)
		if !re.MatchString(scanner.Text()) {
			return &LocationList{}, &LocationList{}, fmt.Errorf("%w: must be two numbers separated by at least one space on each line", ErrInvalidInputFormat)
		}
		values := re.FindStringSubmatch(scanner.Text())

		leftValue, _ := strconv.Atoi(values[1])
		leftLocation, err := NewLocation(leftValue)
		if err != nil {
			return &LocationList{}, &LocationList{}, err
		}

		rightValue, _ := strconv.Atoi(values[2])
		rightLocation, err := NewLocation(rightValue)
		if err != nil {
			return &LocationList{}, &LocationList{}, err
		}

		leftList.AddLocation(leftLocation)
		rightList.AddLocation(rightLocation)
	}

	return leftList, rightList, nil
}

func sumDistances(leftList, rightList *LocationList) int {
	var sum int
	for range leftList.Size() {
		left := leftList.Next()
		right := rightList.Next()
		sum += left.Distance(right)
	}

	return sum
}

func sumSimilarities(leftList, rightList *LocationList) int {
	var sum int
	for range leftList.Size() {
		left := leftList.Next()
		leftMatches := rightList.CountMatches(left)
		similarityScore := left.Int() * leftMatches
		sum += similarityScore
	}

	return sum
}
