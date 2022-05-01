package main

import (
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"

	"google.golang.org/protobuf/proto"
	pb "github.com/sfortson/fitness-tracker/server/proto"
)

func write_proto(data [][]string) {
	var bumps []*pb.Bump
	bf_dict := make(map[string]*pb.Bump)
	for i, line := range data {
		if i == 0 {
			continue
		}
		bf_min, _ := strconv.ParseFloat(line[3], 32)
		bf_max, _ := strconv.ParseFloat(line[4], 32)

		bf_min32 := float32(bf_min)
		bf_max32 := float32(bf_max)

		bfr := &pb.Bump_BodyFatRange{
			Description: &(line[5]),
			HealthRisk: &(line[6]),
			Min:         &bf_min32,
			Max:         &bf_max32,
		}

		min, _ := strconv.ParseInt(line[1], 10, 64)
		max, _ := strconv.ParseInt(line[2], 10, 64)

		age_range := &pb.Bump_AgeRange{
			Min: &min,
			Max: &max,
		}

		if bf_dict[line[1]] != nil {
			bf_dict[line[1]] = &pb.Bump{
				Age:               age_range,
				BodyFatPercentage: append(bf_dict[line[1]].BodyFatPercentage, bfr),
			}
		} else {
			bf_dict[line[1]] = &pb.Bump{
				Age:               age_range,
				BodyFatPercentage: []*pb.Bump_BodyFatRange{bfr},
			}
		}

	}

	for _, j := range bf_dict {
		bumps = append(bumps, j)
	}

	out, err := proto.Marshal(&pb.Bumps{Bump: bumps})

	if err != nil {
		fmt.Println(err)
	}

	if err := ioutil.WriteFile("/Users/sfortson/github-projects/fitness-tracker/server/test.proto", out, 0644); err != nil {
		log.Fatalln("Failed to write address book:", err)
	}
}

func read_csv() {
	// open file
	f, err := os.Open("/Users/sfortson/github-projects/fitness-tracker/server/static/body_fat_table-men.csv")
	if err != nil {
		log.Fatal(err)
	}

	// remember to close the file at the end of the program
	defer f.Close()

	// read csv values using csv.Reader
	csvReader := csv.NewReader(f)
	data, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}
	write_proto(data)
}

func main() {
	read_csv()
}
