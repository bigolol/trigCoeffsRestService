package main

import (
	"encoding/json"
	"net/http"

	"fmt"

	"github.com/bigolol/createTrigCoeffs"
)

type trigRequest struct {
	Amt       int32     `json:"amt"`
	Offset    int32     `json:"offset"`
	AmtPoints int       `json:"amtPoints"`
	PointsX   []float64 `json:"pointsX"`
	PointsY   []float64 `json:"pointsY"`
}

func trigCoeffs(w http.ResponseWriter, r *http.Request) {
	fmt.Println("receiving request")
	var trigReq trigRequest
	//bs := make([]byte, 1024, 1024)
	//amt, _ := r.Body.Read(bs)
	//s := string(bs[:amt])
	//fmt.Println(s)
	err := json.NewDecoder(r.Body).Decode(&trigReq)
	if err != nil {
		panic(err)
	}
	fmt.Println("parsed json")
	defer r.Body.Close()
	linInterpolFunc := func(x float64) float64 {
		i := 0
		fmt.Println(trigReq.AmtPoints)
		for ; i < trigReq.AmtPoints-1 && trigReq.PointsX[i] < x; i++ {
		}
		if trigReq.PointsX[i] == x {
			return trigReq.PointsY[i]
		}
		distFromLast := x - trigReq.PointsX[i-1]
		slope := (trigReq.PointsY[i] - trigReq.PointsY[i-1]) / (trigReq.PointsX[i] - trigReq.PointsX[i-1])
		return trigReq.PointsY[i-1] + slope*distFromLast
	}
	coeffResponse := createTrigCoeffs.CreateCoeffs(linInterpolFunc, trigReq.Amt, trigReq.Offset)
	fmt.Println(coeffResponse.Ais)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(coeffResponse); err != nil {
		panic(err)
	}
}

func main() {
	fmt.Println("starting server")
	http.Handle("/", http.FileServer(http.Dir("./static/")))
	http.HandleFunc("/trigCoeffs", trigCoeffs)
	http.ListenAndServe(":8080", nil)
}
