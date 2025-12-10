package main

import (
	"encoding/json" // แปลง struct ↔ JSON
	"net/http"      // สร้าง API / Server
	"strconv"       // แปลง string ↔ int
	"strings"       // จัดการ string (เช่น split path)
)

// =========================
// 🟢 1. โครงสร้างข้อมูล (Model)
// =========================

type Gallery struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Image  string `json:"image"`
	Detail string `json:"detail"`
}

// =========================
// 🟢 2. Mock Database
// =========================

var galleries = []Gallery{
	{
		ID:     1,
		Name:   "Mochizuki Honami",
		Image:  "/images/Honami_wedding.png", // path ไปยังไฟล์รูปจริง
		Detail: "Mochizuki Honami Wedding Dress Ver.",
	},
	{
		ID:     2,
		Name:   "RX-78-2 Gundam",
		Image:  "/images/gundam.png",
		Detail: "HG 1/144",
	},
	{
		ID:     3,
		Name:   "Usio Noa",
		Image:  "/images/Usio_Noa_Nendoroid.jpg",
		Detail: "Nendoroid Usio Noa",
	},
}

// =========================
// 🟢 3. CORS Middleware
// =========================

func enableCORS(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*") // อนุญาตทุกเว็บเรียก
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

// =========================
// 🟢 4. GET ทั้งหมด + POST เพิ่มข้อมูล
// =========================

func galleriesHandler(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)

	// ✅ ถ้า Browser ส่ง OPTIONS มาก่อน (preflight)
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	// ✅ GET → ดึงข้อมูลทั้งหมด
	if r.Method == "GET" {
		json.NewEncoder(w).Encode(galleries)
		return
	}

	// ✅ POST → เพิ่มข้อมูลใหม่
	if r.Method == "POST" {
		var newGallery Gallery

		// แปลง JSON จาก body → struct
		err := json.NewDecoder(r.Body).Decode(&newGallery)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{
				"message": "Invalid JSON",
			})
			return
		}

		// ✅ สร้าง ID ใหม่อัตโนมัติ
		newGallery.ID = len(galleries) + 1

		// ✅ เพิ่มเข้า mock database
		galleries = append(galleries, newGallery)

		// ✅ ส่ง response กลับไป
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Gallery added successfully",
			"data":    newGallery,
		})
		return
	}

	// ✅ ถ้าไม่ใช่ GET / POST
	w.WriteHeader(http.StatusMethodNotAllowed)
}

// =========================
// 🟢 5. GET ตาม ID
// =========================

func galleryByIDHandler(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	// ✅ ดึง id จาก URL เช่น /api/gallery/2
	path := strings.TrimPrefix(r.URL.Path, "/api/gallery/")
	idStr := path

	// ✅ ถ้าไม่มี id → error
	if idStr == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "ID is required",
		})
		return
	}

	// ✅ แปลง string → int
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Invalid ID",
		})
		return
	}

	// ✅ วนหา gallery ตาม id
	for _, item := range galleries {
		if item.ID == id {
			json.NewEncoder(w).Encode(item)
			return
		}
	}

	// ✅ ถ้าไม่เจอ id
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Gallery not found",
	})
}

// =========================
// 🟢 6. main
// =========================

func main() {

	// ✅ เสิร์ฟไฟล์รูปจากโฟลเดอร์ images
	http.Handle("/images/",
		http.StripPrefix("/images/",
			http.FileServer(http.Dir("images")),
		),
	)

	// ✅ GET ทั้งหมด + POST เพิ่ม
	http.HandleFunc("/api/gallery", galleriesHandler)

	// ✅ GET ตาม id
	http.HandleFunc("/api/gallery/", galleryByIDHandler)
	println("✅ Server running at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
