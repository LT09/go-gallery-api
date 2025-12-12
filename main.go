package main

import (
	"encoding/json" // แปลง struct ↔ JSON
	"net/http"      // ใช้ทำ API / Server
	"strconv"       // แปลง string ↔ int
	"strings"       // ตัด/จัดการ string
)

//
// =========================
// 🟢 1. โครงสร้างข้อมูล (Model)
// =========================
//

type Gallery struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Image  string `json:"image"`
	Detail string `json:"detail"`
}

//
// =========================
// 🟢 2. Mock Database (ข้อมูลชั่วคราวใน slice)
// =========================
//

var galleries = []Gallery{
	{
		ID:     1,
		Name:   "Mochizuki Honami",
		Image:  "/images/Honami_wedding.png",
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

//
// =========================
// 🟢 3. ฟังก์ชันเปิด CORS (ใช้ทุก API)
// =========================
//

func enableCORS(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*") // ทุกเว็บเรียก API ได้
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

//
// =========================
// 🟢 4. Handler หลัก: GET / POST / PUT / DELETE
// =========================
//

func galleryHandler(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)

	// Preflight → OPTIONS (Browser จะยิงมาก่อน POST/PUT/DELETE)
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	// ทุก response เป็น JSON
	w.Header().Set("Content-Type", "application/json")

	//
	// 🟡 ดึง ID จาก URL
	// เช่น /api/gallery/5 → "5"
	//
	path := strings.TrimPrefix(r.URL.Path, "/api/gallery") // แก้จุด error หลัก
	path = strings.TrimPrefix(path, "/")                    // ตัด "/" หน้าออก ถ้ามี
	idStr := path                                           // ถ้าว่าง = ไม่มี ID

	//
	// 🟢 GET ทั้งหมด (ไม่มี ID)
	//
	if r.Method == "GET" && idStr == "" {
		json.NewEncoder(w).Encode(galleries)
		return
	}

	//
	// 🟡 มี ID → แปลงเป็น int
	//
	var id int
	var err error
	if idStr != "" {
		id, err = strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid ID", http.StatusBadRequest)
			return
		}
	}

	//
	// 🟢 GET by ID
	//
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

	//
	// 🟢 POST → เพิ่มข้อมูลใหม่
	//
	if r.Method == "POST" {

		var newItem Gallery

		// อ่าน JSON body → struct
		err := json.NewDecoder(r.Body).Decode(&newItem)
		if err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		// Auto ID (ง่ายๆ)
		newItem.ID = len(galleries) + 1

		galleries = append(galleries, newItem)

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(newItem)
		return
	}

	//
	// 🟢 PUT แก้ไขข้อมูลตาม ID
	//
	if r.Method == "PUT" {
		var updateItem Gallery

		err := json.NewDecoder(r.Body).Decode(&updateItem)
		if err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		for i, g := range galleries {
			if g.ID == id {

				// อัปเดตข้อมูล
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

	//
	// 🟢 DELETE - ลบข้อมูลตาม ID
	//
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

	//
	// ❌ Method ที่ไม่รองรับ
	//
	http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
}

//
// =========================
// 🟢 5. main function
// =========================
//

func main() {

	// เสิร์ฟไฟล์รูปจริง
	// /images/... → ไปหยิบจากโฟลเดอร์ images
	http.Handle("/images/",
		http.StripPrefix("/images/",
			http.FileServer(http.Dir("images")),
		),
	)

	// ❗️❗️ สำคัญมาก: ต้องให้รองรับทั้ง /api/gallery และ /api/gallery/
	// ไม่งั้น POST /api/gallery จะไม่เข้า handler
	http.HandleFunc("/api/gallery", galleryHandler)  // แบบไม่มี slash ท้าย
	http.HandleFunc("/api/gallery/", galleryHandler) // แบบมี slash ท้าย

	println("🚀 Server running at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
