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
var mock bool
var threads = make(chan [2]string, 100)
var root string
var personDir string
var companyDir string
var audioDir string
var coverNum int
var activeSupplier []int

func loadConfig(){
    config := lib.NewConfig()

    debug = config.Debug
    mock = config.Mock

    root = config.Resource.Root
    personDir  = config.Resource.Person
    companyDir  = config.Resource.Company
    audioDir  = config.Resource.AudioPath
    coverNum = config.Resource.CoverNum
}

/**
 * 供应商文件同步
 *
 *      按供应商类型目录进行同步。
 *      1. 检查到供应商时，会分出一个进程来检查该供应商的作品类型；
 *          level2-> 检查到作品类型时，会分出一个进程来检查该作品类型下的文件；
 *              level3-> 检查到文件时，如果文件没有被转码，且文件是一天内创建的，会加入文件转码队列。
 *      2. 供应商检查完毕时，会开启文件转码过程，
 *      3. 当开始转码过程时，level2与level3可能会依然在执行，如果文件转码队列为空，且在wait时长内，没有检查到新文件，会关闭程序，
 *      然而这时level2或level3的进程还没有运行完毕（旧文件太多，遍历时间过长），导致同步数据不完整，所以要注意适当加大wait的值
 * 
 * @method  CountSupplier
 * @author  万引  v.songzp@yoozoo.com  2019-05-21T10:54:50+0800
 */
func CountSupplier(){
    loadConfig();

    // 初始化任务控制 channel
    service.InitQueue()

    var personPath = root + personDir
    countPerson(personPath)

    var companyPath = root + companyDir
    countCompany(companyPath)

    var audioPath = root + audioDir
    countAudio(audioPath)

    // TODO 更新activeSupplier以外的项目为已删除
    if !mock {
        DeleteSupplierExcept(activeSupplier)    
    }
    

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

func countAudio(path string){
    var supplierType = &SupplierType{3, audioDir}
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
        //  fmt.Println("[DiscoverCategory]:", categoryName)
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
    if !mock {
        DeleteCategory(supplierId, categoryIds)    
    }
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
        //  fmt.Println("[DiscoverFile]:", fileName)
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
    if !mock {
        DeleteFile(categoryId, fileIds)    
    }
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
func GetSupplierCovers(supplierUrl string, supplierType int) (covers []string, err error) {
    var realPath = root + supplierUrl
    var coverPath []string


    if supplierType == 3 {
        coverPath = []string{"\\folder\\folder_0.png","\\folder\\folder_0.png","\\folder\\folder_0.png"}
    }else{

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

