Example of adding a new type Time `time.Time`. 

```
./fields/
├── comparators
│   └── time.go     -- Comparator for make query by field with type time.Time
├── getter.go       -- Getter for a time.Time type
├── indexes
│   └── time.go     -- Optional index for optimize querying by field with type time.Time  
└── querybuilder
    ├── base.go     -- Extend query.BaseQueryBuilder with method WhereTime()
    └── debug.go    -- Save for debug & logging query by field with type time.Time 
```
