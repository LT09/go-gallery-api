package main // บอก Go ว่านี่คือโปรแกรมหลักที่รันได้จริง

import (
	"encoding/json" // ใช้แปลง struct → JSON
	"net/http"      // ใช้สร้าง web server และ API
	"strconv"       // ใช้แปลง string → int
)

// =========================
// 🟢 1. โครงสร้างข้อมูล (Model)
// =========================

// Gallery คือโครงสร้างข้อมูลของรูป 1 ชิ้น
// เทียบได้กับ interface หรือ type ใน TypeScript
type Gallery struct {
	ID     int    `json:"id"`     // id ของรูป
	Name   string `json:"name"`   // ชื่อรูป
	Image  string `json:"image"`  // path ของไฟล์รูป
	Detail string `json:"detail"` // รายละเอียด
}


// =========================
// 🟢 2. Mock Database (จำลองฐานข้อมูล)
// =========================

// galleries คือข้อมูลจำลองที่ใช้แทนฐานข้อมูลจริง
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
// 🟢 3. API Handler (Controller)
// =========================

// galleryHandler จะทำงานเมื่อมีคนเรียก /api/gallery
func galleryHandler(w http.ResponseWriter, r *http.Request) {

	// ✅ CORS: อนุญาตให้ทุกเว็บเรียก API นี้ได้ (ใช้ตอนพัฒนา)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// ✅ ถ้าเป็น preflight (OPTIONS) ให้จบทันที
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	// ✅ บอก client ว่าจะส่ง JSON กลับไป
	w.Header().Set("Content-Type", "application/json")

	// ✅ ดึงค่า id จาก query string เช่น ?id=2
	idStr := r.URL.Query().Get("id")

	// ✅ ถ้าไม่มี id → ส่งข้อมูลทั้งหมด
	if idStr == "" {
		// ✅ แปลง struct → JSON และส่งกลับไปทันที
	json.NewEncoder(w).Encode(galleries)
		return
	}

	// ✅ แปลง id จาก string → int
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Invalid ID",
		})
		return
	}

	// ✅ วน loop หา gallery ที่ id ตรงกัน
	for _, g := range galleries {
		if g.ID == id {
			json.NewEncoder(w).Encode(g)
			return
		}
	}

	// ✅ ถ้า loop ครบแล้วไม่เจอ
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Gallery not found",
	})
}


// =========================
// 🟢 4. main() จุดเริ่มต้นโปรแกรม
// =========================

func main() {

	// ✅ ให้ Go Server เสิร์ฟไฟล์ทุกไฟล์ในโฟลเดอร์ images
	// ถ้า browser เรียก /images/xxx.jpg → ไปอ่านไฟล์จากโฟลเดอร์ images
	http.Handle(
		"/images/",
		http.StripPrefix(
			"/images/",
			http.FileServer(http.Dir("images")),
		),
	)

	// ✅ ผูก API /api/gallery กับฟังก์ชัน galleryHandler
	http.HandleFunc("/api/gallery", galleryHandler)

	// ✅ แสดงข้อความใน terminal
	println("✅ Server running at http://localhost:8080")

	// ✅ เปิด server ที่ port 8080 และรอ request ไปเรื่อยๆ
	http.ListenAndServe(":8080", nil)
}
