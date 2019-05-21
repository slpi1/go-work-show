package resource

import (
    "encoding/json"

    "lib"
    "resource/model"
    "service"
)

func SaveSupplier(name string, supplierType *SupplierType)(supplierId int, err error) {
    var supplier = &model.Supplier{}
    var engine = lib.Connection()
    
    if !mock {
        if has,err := engine.Where("name = ?", name).Get(supplier); !has {
            if err != nil {
                return 0, err;
            }
        }
    }

    supplier.Name = name
    supplier.Type = supplierType.Type
    supplier.Url = "\\" + supplierType.Prefix + "\\" + name
    supplier.IsDelete = 0

    covers, err := GetSupplierCovers(supplierType.Prefix + "\\" + name)
    if err != nil {
        lib.Logger().Println(supplier.Url, "covers get failed")
    }

    getSupplierAttribute(covers, supplier)

    if mock {
        return 1,nil
    }

    if supplier.Id > 0 {
        _, err := engine.ID(supplier.Id).AllCols().Update(supplier)
        if err != nil {
            lib.Logger().Println("update failed")
            return 0, err
        }
    } else {
        _, err := engine.InsertOne(supplier)
        if err != nil {
            lib.Logger().Println("insert failed")
            return 0, err
        }
    }
    
    supplierId = supplier.Id
    return supplierId, nil
}

func GetCoverInfo(covers []string)(coverPath, coverCompressPath1, coverCompressPath2 []string){

    for _,file := range covers {
        compress1 := service.GetThumbPath(file)
        compress2 := service.GetPreviewPath(file)

        coverPath = append(coverPath, file)
        coverCompressPath1 = append(coverCompressPath1, compress1)
        coverCompressPath2 = append(coverCompressPath2, compress2)
    }
    return coverPath,coverCompressPath1,coverCompressPath2
}

func getSupplierAttribute(covers []string, supplier *model.Supplier) error {
    //作品类型封面路径
    coverPath,coverCompressPath1,coverCompressPath2 := GetCoverInfo(covers)

    jsonData, err := json.Marshal(coverPath)
    if err != nil {
        lib.Logger().Println("type json error", err.Error())
        return err
    }
    supplier.CoverPath = string(jsonData)

    jsonData, err = json.Marshal(coverCompressPath1)
    if err != nil {
        lib.Logger().Println("type json error", err.Error())
        return err
    }
    supplier.CoverCompressPath1 = string(jsonData)

    jsonData, err = json.Marshal(coverCompressPath2)
    if err != nil {
        lib.Logger().Println("type json error", err.Error())
        return err
    }
    supplier.CoverCompressPath2 = string(jsonData)
    return nil
}

func DeleteSupplierExcept( activeSupplierIds []int) {

    var engine = lib.Connection()

    supplier := new(model.Supplier)
    supplier.IsDelete = 1

    _, err := engine.NotIn("id", activeSupplierIds).Update(supplier)
    if err != nil {
        lib.Logger().Println("Delete supplier Failed", err.Error())
    }
}