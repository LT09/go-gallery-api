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
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

// =========================
// 🟢 4. GET ทั้งหมด + POST เพิ่มข้อมูล
// =========================

func galleryHandler(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)

	// Preflight → OPTIONS
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	// ดึง path หลัง /api/gallery/
	path := strings.TrimPrefix(r.URL.Path, "/api/gallery/")
	idStr := path // ถ้าว่าง = ไม่มี ID

	// -------------------------
	// GET ทั้งหมด
	// -------------------------
	if r.Method == "GET" && idStr == "" {
		json.NewEncoder(w).Encode(galleries)
		return
	}

	// -------------------------
	// มี ID → ต้องแปลงเป็น int
	// -------------------------
	var id int
	var err error
	if idStr != "" {
		id, err = strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid ID", http.StatusBadRequest)
			return
		}
	}

	// -------------------------
	// GET by ID
	// -------------------------
	if r.Method == "GET" {
		for _, g := range galleries {
			if g.ID == id {
				json.NewEncoder(w).Encode(g)
				return
			}
		}
		http.Error(w, "Gallery not found", http.StatusNotFound)
		return
	}

	// -------------------------
	// POST - เพิ่มข้อมูลใหม่
	// -------------------------
	if r.Method == "POST" {
		var newItem Gallery

		// decode JSON body → struct
		err := json.NewDecoder(r.Body).Decode(&newItem)
		if err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		// Auto ID
		newItem.ID = len(galleries) + 1
		galleries = append(galleries, newItem)

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(newItem)
		return
	}

	// -------------------------
	// PUT - แก้ไขข้อมูลตาม ID
	// -------------------------
	if r.Method == "PUT" {
		var updateItem Gallery

		err := json.NewDecoder(r.Body).Decode(&updateItem)
		if err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		for i, g := range galleries {
			if g.ID == id {

				// อัปเดตข้อมูลใหม่
				galleries[i].Name = updateItem.Name
				galleries[i].Image = updateItem.Image
				galleries[i].Detail = updateItem.Detail

				json.NewEncoder(w).Encode(galleries[i])
				return
			}
		}

		http.Error(w, "Gallery not found", http.StatusNotFound)
		return
	}

	// -------------------------
	// DELETE - ลบข้อมูลตาม ID
	// -------------------------
	if r.Method == "DELETE" {
		for i, g := range galleries {
			if g.ID == id {

				// ลบ index i ออกจาก slice
				galleries = append(galleries[:i], galleries[i+1:]...)

				w.WriteHeader(http.StatusOK)
				w.Write([]byte("Deleted successfully"))
				return
			}
		}

		http.Error(w, "Gallery not found", http.StatusNotFound)
		return
	}

	// Method ไม่รองรับ
	http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
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

	// ✅ GET POST PUT DELETE /api/gallery/
	http.HandleFunc("/api/gallery/", galleryHandler)
	println("✅ Server running at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
