package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"api-students/app/model"
	"api-students/config"
	"api-students/database"
	"api-students/helper"
)

var (
	baseURL   = "http://localhost:3000/api/v1"
	jwtSecret string
)

func main() {
	mode := flag.String("mode", "crowd", "Mode uji beban: crowd | credits | swap | supervisors")
	flag.Parse()

	cfg := config.Load()
	jwtSecret = cfg.JWTSecret
	if jwtSecret == "" {
		jwtSecret = "65b458dc012be6c16665db71ee6ec7f53dc9726237bbbd4e7288bacd789955cf"
	}

	jwtManager := helper.NewJWTManager(jwtSecret, "praktikum-backend", 15*time.Minute)

	tokens := make([]string, 61)
	userIDs := make([]int, 61)

	db, err := database.NewDBPool(cfg)
	if err == nil {
		defer db.Close()
		rows, qErr := db.Query(context.Background(), "SELECT id, username, email, role FROM users WHERE username LIKE 'mhs%' ORDER BY username")
		if qErr == nil {
			defer rows.Close()
			idx := 1
			for rows.Next() && idx <= 60 {
				var u model.AuthUser
				if err := rows.Scan(&u.UserID, &u.Username, &u.Email, &u.Role); err == nil {
					userIDs[idx] = u.UserID
					token, tErr := jwtManager.GenerateAccessToken(u)
					if tErr == nil {
						tokens[idx] = token
					}
					idx++
				}
			}
		}
	}

	// Fallback if db query failed
	for i := 1; i <= 60; i++ {
		if tokens[i] == "" {
			username := fmt.Sprintf("mhs%03d", i)
			u := model.AuthUser{
				UserID:   userIDs[i],
				Username: username,
				Email:    fmt.Sprintf("%s@student.unair.ac.id", username),
				Role:     "user",
			}
			token, _ := jwtManager.GenerateAccessToken(u)
			tokens[i] = token
		}
	}

	start := time.Now()

	switch *mode {
	case "crowd":
		runCrowd(tokens)
	case "credits":
		runCredits(tokens)
	case "swap":
		runSwap(tokens)
	case "supervisors":
		runSupervisors(tokens)
	default:
		log.Fatalf("mode tidak dikenal: %s", *mode)
	}

	elapsed := time.Since(start)
	fmt.Printf("waktu total : %v\n", elapsed)
}

// fire menjalankan n pekerjaan secara serentak dan menghitung hasilnya.
// Bug #1 fixed: out chan<- string (write-only channel)
func fire(n int, job func(i int) string) map[string]int {
	var wg sync.WaitGroup
	start := make(chan struct{})
	results := make(chan string, n)

	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(i int, out chan<- string) {
			defer wg.Done()
			<-start
			out <- job(i)
		}(i, results)
	}

	close(start) // buka gerbang serentak
	wg.Wait()
	close(results)

	tally := make(map[string]int)
	for r := range results {
		tally[r]++
	}
	return tally
}

func getCourseID(code string, token string) int {
	req, _ := http.NewRequest("GET", baseURL+"/courses", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0
	}
	defer resp.Body.Close()

	var res struct {
		Data []struct {
			ID   int    `json:"id"`
			Code string `json:"code"`
		} `json:"data"`
	}
	_ = json.NewDecoder(resp.Body).Decode(&res)
	for _, c := range res.Data {
		if c.Code == code {
			return c.ID
		}
	}
	return 0
}

func runCrowd(tokens []string) {
	fmt.Println("60 mahasiswa berebut kelas IF399 secara serentak...")
	courseID := getCourseID("IF399", tokens[1])
	if courseID == 0 {
		courseID = 13 // fallback
	}

	tally := fire(60, func(i int) string {
		mhsIdx := i + 1
		token := tokens[mhsIdx]

		req, err := http.NewRequest("POST", fmt.Sprintf("%s/courses/%d/enroll", baseURL, courseID), bytes.NewBuffer([]byte("{}")))
		if err != nil {
			return "ERR"
		}
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Content-Type", "application/json")
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return "ERR"
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusCreated {
			return "201 OK"
		}

		body, _ := io.ReadAll(resp.Body)
		var errResp struct {
			Code string `json:"code"`
		}
		_ = json.Unmarshal(body, &errResp)
		if errResp.Code != "" {
			return fmt.Sprintf("%d %s", resp.StatusCode, errResp.Code)
		}
		return fmt.Sprintf("%d", resp.StatusCode)
	})

	fmt.Println("hasil request:")
	for status, count := range tally {
		fmt.Printf("  %-25s %d\n", status, count)
	}
}

func runCredits(tokens []string) {
	fmt.Println("10 mahasiswa, masing-masing mendaftar ke 12 kelas 3 SKS secara serentak (batas 24 SKS)...")
	// Ambil daftar kelas
	req, _ := http.NewRequest("GET", baseURL+"/courses", nil)
	req.Header.Set("Authorization", "Bearer "+tokens[1])
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	var res struct {
		Data []struct {
			ID   int    `json:"id"`
			Code string `json:"code"`
		} `json:"data"`
	}
	_ = json.NewDecoder(resp.Body).Decode(&res)

	var targetCourseIDs []int
	for _, c := range res.Data {
		if c.Code != "IF399" && len(targetCourseIDs) < 12 {
			targetCourseIDs = append(targetCourseIDs, c.ID)
		}
	}

	totalJobs := 10 * len(targetCourseIDs)
	tally := fire(totalJobs, func(idx int) string {
		studentNum := (idx / len(targetCourseIDs)) + 1
		courseIdx := idx % len(targetCourseIDs)
		cID := targetCourseIDs[courseIdx]
		token := tokens[studentNum]

		postReq, err := http.NewRequest("POST", fmt.Sprintf("%s/courses/%d/enroll", baseURL, cID), bytes.NewBuffer([]byte("{}")))
		if err != nil {
			return "ERR"
		}
		postReq.Header.Set("Authorization", "Bearer "+token)
		postReq.Header.Set("Content-Type", "application/json")
		resp, err := http.DefaultClient.Do(postReq)
		if err != nil {
			return "ERR"
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusCreated {
			return "201 OK"
		}

		body, _ := io.ReadAll(resp.Body)
		var errResp struct {
			Code string `json:"code"`
		}
		_ = json.Unmarshal(body, &errResp)
		if errResp.Code != "" {
			return fmt.Sprintf("%d %s", resp.StatusCode, errResp.Code)
		}
		return fmt.Sprintf("%d", resp.StatusCode)
	})

	fmt.Println("hasil request:")
	for status, count := range tally {
		fmt.Printf("  %-25s %d\n", status, count)
	}
}

func runSwap(tokens []string) {
	fmt.Println("20 pasang mahasiswa pindah kelas berlawanan arah secara serentak...")
	c1ID := getCourseID("IF101", tokens[1])
	c2ID := getCourseID("IF102", tokens[1])

	// Pastikan 20 mahasiswa pertama terdaftar di IF101 dan 20 berikutnya terdaftar di IF102
	for i := 1; i <= 20; i++ {
		req, _ := http.NewRequest("POST", fmt.Sprintf("%s/courses/%d/enroll", baseURL, c1ID), bytes.NewBuffer([]byte("{}")))
		req.Header.Set("Authorization", "Bearer "+tokens[i])
		req.Header.Set("Content-Type", "application/json")
		resp, err := http.DefaultClient.Do(req)
		if err == nil {
			resp.Body.Close()
		}
	}
	for i := 21; i <= 40; i++ {
		req, _ := http.NewRequest("POST", fmt.Sprintf("%s/courses/%d/enroll", baseURL, c2ID), bytes.NewBuffer([]byte("{}")))
		req.Header.Set("Authorization", "Bearer "+tokens[i])
		req.Header.Set("Content-Type", "application/json")
		resp, err := http.DefaultClient.Do(req)
		if err == nil {
			resp.Body.Close()
		}
	}

	// Sekarang 40 mahasiswa swap serentak
	tally := fire(40, func(idx int) string {
		studentNum := idx + 1
		token := tokens[studentNum]

		var payload []byte
		if studentNum <= 20 {
			payload, _ = json.Marshal(map[string]int{"from_course_id": c1ID, "to_course_id": c2ID})
		} else {
			payload, _ = json.Marshal(map[string]int{"from_course_id": c2ID, "to_course_id": c1ID})
		}

		req, err := http.NewRequest("POST", baseURL+"/enrollments/swap", bytes.NewReader(payload))
		if err != nil {
			return "ERR"
		}
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Content-Type", "application/json")
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return "ERR"
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			return "200 OK"
		}

		body, _ := io.ReadAll(resp.Body)
		var errResp struct {
			Code string `json:"code"`
		}
		_ = json.Unmarshal(body, &errResp)
		if errResp.Code != "" {
			return fmt.Sprintf("%d %s", resp.StatusCode, errResp.Code)
		}
		return fmt.Sprintf("%d", resp.StatusCode)
	})

	fmt.Println("hasil request:")
	for status, count := range tally {
		fmt.Printf("  %-25s %d\n", status, count)
	}
}

func runSupervisors(tokens []string) {
	fmt.Println("30 penetapan serentak ke dosen pembimbing berkuota 3...")
	adminUser := model.AuthUser{
		UserID:   9999,
		Username: "admin",
		Email:    "admin@unair.ac.id",
		Role:     "admin",
	}
	jwtManager := helper.NewJWTManager(jwtSecret, "praktikum-backend", 15*time.Minute)
	adminToken, _ := jwtManager.GenerateAccessToken(adminUser)

	studentIDs := make([]int, 0, 30)
	cfg := config.Load()
	db, err := database.NewDBPool(cfg)
	if err == nil {
		defer db.Close()
		rows, _ := db.Query(context.Background(), "SELECT id FROM students WHERE deleted_at IS NULL ORDER BY id LIMIT 30")
		if rows != nil {
			defer rows.Close()
			for rows.Next() {
				var sID int
				if err := rows.Scan(&sID); err == nil {
					studentIDs = append(studentIDs, sID)
				}
			}
		}
	}

	tally := fire(len(studentIDs), func(idx int) string {
		studentID := studentIDs[idx]
		payload, _ := json.Marshal(map[string]int{"student_id": studentID})

		req, err := http.NewRequest("POST", baseURL+"/supervisors/1/assign", bytes.NewReader(payload))
		if err != nil {
			return "ERR"
		}
		req.Header.Set("Authorization", "Bearer "+adminToken)
		req.Header.Set("Content-Type", "application/json")
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return "ERR"
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusCreated {
			return "201 OK"
		}

		body, _ := io.ReadAll(resp.Body)
		var errResp struct {
			Code string `json:"code"`
		}
		_ = json.Unmarshal(body, &errResp)
		if errResp.Code != "" {
			return fmt.Sprintf("%d %s", resp.StatusCode, errResp.Code)
		}
		return strconv.Itoa(resp.StatusCode)
	})

	fmt.Println("hasil request:")
	for status, count := range tally {
		fmt.Printf("  %-25s %d\n", status, count)
	}
}
