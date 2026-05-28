package Employees

func (UC *UserControl) AddSubscriptionRoutes(R *Routes) {
	var r Route
	var Module, Level1 string
	var ModuleDisplayOrder, Level1DisplayOrder int64
	Module = "Value Added Services" //Configuration.Module
	ModuleDisplayOrder = 1
	Level1 = "Product Design Center"
	Level1DisplayOrder = 1

	Level1 = "Employees"
	Level1DisplayOrder = Level1DisplayOrder + 1
	r = Route{
		"HTTP_Employees",
		"GET",
		"/" + Configuration.Module + "/" + Configuration.Version + "/HTTP_Employees/",
		Use(UC.HTTP_Employees),
		true,
		"Employees - Read", // DisplayName
		1,                  // DisplayOrder
		Module,             // Module
		ModuleDisplayOrder, //ModuleDisplayOrder
		Level1,             // Level1
		Level1DisplayOrder, // Level1DisplayOrder
		"",                 // Level2
		0,                  // Level2DisplayOrder
		"",                 // Level3
		0,                  // Level3DisplayOrder
	}
	*R = append(*R, r)

	r = Route{
		"HTTP_Employees",
		"POST",
		"/" + Configuration.Module + "/" + Configuration.Version + "/HTTP_Employees/",
		Use(UC.HTTP_Employees),
		true,
		"Employees - Add",  // DisplayName
		1,                  // DisplayOrder
		Module,             // Module
		ModuleDisplayOrder, //ModuleDisplayOrder
		Level1,             // Level1
		Level1DisplayOrder, // Level1DisplayOrder
		"",                 // Level2
		0,                  // Level2DisplayOrder
		"",                 // Level3
		0,                  // Level3DisplayOrder
	}
	*R = append(*R, r)

	r = Route{
		"HTTP_Employees",
		"PUT",
		"/" + Configuration.Module + "/" + Configuration.Version + "/HTTP_Employees/{Login}",
		Use(UC.HTTP_Employees),
		true,
		"Employees - Edit", // DisplayName
		1,                  // DisplayOrder
		Module,             // Module
		ModuleDisplayOrder, //ModuleDisplayOrder
		Level1,             // Level1
		Level1DisplayOrder, // Level1DisplayOrder
		"",                 // Level2
		0,                  // Level2DisplayOrder
		"",                 // Level3
		0,                  // Level3DisplayOrder
	}
	*R = append(*R, r)

	r = Route{
		"HTTP_Employees",
		"DELETE",
		"/" + Configuration.Module + "/" + Configuration.Version + "/HTTP_Employees/{Login}",
		Use(UC.HTTP_Employees),
		true,
		"Employees - Delete", // DisplayName
		1,                    // DisplayOrder
		Module,               // Module
		ModuleDisplayOrder,   //ModuleDisplayOrder
		Level1,               // Level1
		Level1DisplayOrder,   // Level1DisplayOrder
		"",                   // Level2
		0,                    // Level2DisplayOrder
		"",                   // Level3
		0,                    // Level3DisplayOrder
	}
	*R = append(*R, r)
}
