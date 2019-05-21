package resource

import (
	"io/ioutil"
    "os"
    "fmt"

    "lib"
    "service"
)

type SupplierType struct {
	Type int
	Prefix string
}

var debug bool
var threads = make(chan [2]string, 100)
var root string
var personDir string
var companyDir string
var coverNum int
var activeSupplier []int

func loadConfig(){
	config := lib.NewConfig()

	debug = config.Debug
	root = config.Resource.Root
	personDir  = config.Resource.Person
	companyDir  = config.Resource.Company
	coverNum = config.Resource.CoverNum
}

// 遍历供应商
func CountSupplier(){
	loadConfig();

	// 初始化任务控制 channel
    service.InitQueue()

	var personPath = root + personDir
	countPerson(personPath)

	var companyPath = root + companyDir
	countCompany(companyPath)

	// TODO 更新activeSupplier以外的项目为已删除
	DeleteSupplierExcept(activeSupplier)

    // 执行转码
    service.StartExec()
}

func countCompany(path string){
	var supplierType = &SupplierType{1, companyDir}
	DiscoverSupplier(path, supplierType)
}

func countPerson(path string){
	var supplierType = &SupplierType{2, personDir}
	DiscoverSupplier(path, supplierType)
}


// 发现供应商
func DiscoverSupplier(path string, supplierType *SupplierType) error {
	supplierDirs := scanDir(path)

	for _, supplierDir := range supplierDirs {
		supplierName :=  supplierDir.Name()
		
		supplierUrl := supplierType.Prefix + "\\" + supplierName
		if debug {
			fmt.Println("[DiscoverSupplier]:", supplierName)
		}

		// 更新供应商数据
		supplierId,err := SaveSupplier(supplierName, supplierType)
		if err != nil {
			lib.Logger().Println("Save Supplier Failed:", supplierUrl)
			continue
		}
		activeSupplier = append(activeSupplier, supplierId)

		go DiscoverCategory(supplierUrl, supplierId)
	}
	return nil
}

// 发现供应商分类
func DiscoverCategory(supplierUrl string, supplierId int) error {
	var realPath = root + supplierUrl
	var categoryIds []int

	categories := scanDir(realPath)

	for _, category := range categories {

		if !category.IsDir() {
			continue
		}

		categoryName :=  category.Name()
		// if debug {
		// 	fmt.Println("[DiscoverCategory]:", categoryName)
		// }
		
		categoryUrl := supplierUrl + "\\" + categoryName

		// 更新作品分类数据
		categoryId,err := SaveProductType(categoryName,supplierId, supplierUrl)
		if err != nil {
			lib.Logger().Println("Save  Supplier Category Failed:", categoryUrl)
			continue
		}
		categoryIds = append(categoryIds, categoryId)

		go DiscoverFile(categoryUrl, categoryId)
	}

	// TODO 更新supplierId下，categoryIds以外的项目为已删除
	DeleteCategory(supplierId, categoryIds)
	return nil
}

// 发现作品
func DiscoverFile(categoryUrl string, categoryId int) error {
	var realPath = root + categoryUrl
	var fileIds []int

	files := scanDir(realPath)

	for _, file := range files {

		if file.IsDir() || !service.CheckFileType(file) {
			continue
		}

		fileName :=  file.Name()
		// if debug {
		// 	fmt.Println("[DiscoverFile]:", fileName)
		// }

		// 更新作品文件数据
		fileId,err := SaveProductFile(fileName,categoryId, categoryUrl)
		if err != nil {
			lib.Logger().Println("Save Supplier Category File Failed:", categoryUrl, fileName)
			continue
		}
		fileIds = append(fileIds, fileId)

		service.FormatFile(realPath + "\\" + fileName)
	}

	// TODO 更新categoryId下，fileIds以外的项目为已删除
	DeleteFile(categoryId, fileIds)
	return nil
}

// 遍历目录
func scanDir(path string) []os.FileInfo {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		lib.Logger().Println("Scan Path Error:", path, err.Error())
		return nil
	}
	return files
}


// 获取供应商封面文件列表
func GetSupplierCovers(supplierUrl string) (covers []string, err error) {
	var realPath = root + supplierUrl
	var coverPath []string

	dirs := scanDir(realPath)

	for _, category := range dirs {
		if len(coverPath) < coverNum {

			if !category.IsDir() {
				continue
			}

			categoryCovers,err := GetCategoryCovers(supplierUrl, category.Name())
			if err != nil {
				continue
			}

			coverPath = append(coverPath, categoryCovers...)
		}
		
	}

	if len(coverPath) > 3 {
		coverPath = coverPath[:3]
	}
	return coverPath, nil
}

// 获取供应商分类的封面文件列表
func GetCategoryCovers(supplierUrl string, categoryName string)(covers []string, err error){
	var realPath = root + supplierUrl + "\\" + categoryName
	var coverPath []string

	dirs := scanDir(realPath)

	for _, file := range dirs {
		if len(coverPath) < coverNum {
			if file.IsDir() || !service.CheckFileType(file) {
				continue
			}

			var filePath = "\\" + supplierUrl + "\\" + categoryName + "\\" + file.Name()
			coverPath = append(coverPath, filePath)
		}
		
	}
	return coverPath, nil
}

