package resource

import (
	"encoding/json"
	"fmt"

	"lib"
    "resource/model"
)


func SaveProductType(name string, supplierId int, supplierUrl string)(categoryId int, err error) {
	var category = &model.ProductType{}

	var engine = lib.Connection()
	if has,err := engine.Where("name = ?", name).And("supplier_id = ?", supplierId).Get(category); !has {
		if err != nil {
			return 0, err;
		}
	}

	category.SupplierId = supplierId
	category.Name = name
	category.Url = "\\" + supplierUrl + "\\" + name
	category.ParentId = -1
	category.IsDelete = 0

	covers, err := GetCategoryCovers(supplierUrl, name)
	if err != nil {
		lib.Logger().Println(category.Url, "covers get failed")
	}

	dirty := getcategoryAttribute(covers, category)

	if category.Id > 0 {
		if dirty {
			if debug {
				fmt.Println("category info update:", name)
			}

			_, err := engine.ID(category.Id).AllCols().Update(category)
			if err != nil {
				lib.Logger().Println("update failed")
				return 0, err
			}
		}
	} else {
		_, err := engine.InsertOne(category)
		if err != nil {
			lib.Logger().Println("insert failed")
			return 0, err
		}
	}
	
	categoryId = category.Id
	return categoryId, nil
}



func getcategoryAttribute(covers []string, category *model.ProductType) bool {
	var dirty = false
	//作品类型封面路径
	coverPath,coverCompressPath1,coverCompressPath2 := GetCoverInfo(covers)

	jsonData, err := json.Marshal(coverPath)
	if err != nil {
		lib.Logger().Println("type json error", err.Error())
		return false
	}
	coverPathString := string(jsonData)
	if category.CoverPath != coverPathString {
		category.CoverPath = coverPathString
		dirty = true
	}

	jsonData, err = json.Marshal(coverCompressPath1)
	if err != nil {
		lib.Logger().Println("type json error", err.Error())
		return false
	}
	coverCompressPath1String := string(jsonData)

	if category.CoverCompressPath1 != coverCompressPath1String {
		category.CoverCompressPath1 = coverCompressPath1String
		dirty = true
	}

	jsonData, err = json.Marshal(coverCompressPath2)
	if err != nil {
		lib.Logger().Println("type json error", err.Error())
		return false
	}
	coverCompressPath2String := string(jsonData)
	if category.CoverCompressPath2 != coverCompressPath2String {
		category.CoverCompressPath2 = coverCompressPath2String
		dirty = true
	}
	return dirty
}



func DeleteCategory(supplierId int, activeCategoryIds []int) {

	var engine = lib.Connection()

	producttype := new(model.ProductType)
	producttype.IsDelete = 1

	_, err := engine.Where("supplier_id = ?", supplierId).NotIn("id", activeCategoryIds).Update(producttype)
	if err != nil {
		lib.Logger().Println("Delete Category Failed", err.Error())
	}
}