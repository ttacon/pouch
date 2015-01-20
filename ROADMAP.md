ROADMAP
=======

 - [ ] Pouch implementations (not in any particular order)
   - [ ] Postgres
   - [ ] sqlite
   - [ ] Redis
   - [ ] Mongo (?)
 - [ ] How to generate different code for multiple pouch types
   - [ ] look into maybe adding another param to every interface (PouchType const or just string?)
   - [ ] Add code generation for each as we go along
 - [ ] Code generation in general
   - [✔] Ability to distinguish differences between what's in db (schema-wise) and struct
   - [✔] Ability to create tables/etc from structs
   - [✔] Generate method implementations from structs
   - [ ] wrapper types for structs in other packages?
   - [ ] Nested object hoisting (?)
 - [ ] MOAR tests
   - [ ] Make moc executors etc
   - [✔] Get to 20% test coverage
   - [ ] Get to 40% test coverage
   - [ ] Get to 60% test coverage
   - [ ] Get to 70% test coverage
   - [ ] Get to 80% test coverage (like this will ever happen)
   - [ ] Get to 90% test coverage
 - [ ] Work on interface design (in reference to different pouch types as well)
 - [ ] Joins
   - [ ] Allow for some kind of laziness?
   - [ ] Generate getters, setters? 
     - [ ] Insert closure? vs
     - [ ] Embed interface (```Retriever```?)
 - [ ] Make pouch tool faster
   - [ ] Profile pouch to see what is taking forever
   - [ ] Make the slow stuff faster

