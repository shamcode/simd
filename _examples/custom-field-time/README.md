Example of adding a new type Time `time.Time`. 

```
./types/
├── comparators
│   └── time.go     -- Comparator for make query by field with type time.Time
├── getter.go       -- Getter for a time.Time type
├── indexes
│   └── time.go     -- Optional index for optimize querying by field with type time.Time  
└── querybuilder
    └── options.go  -- Add new query.BuilderOption for build condition by time.Time
```
