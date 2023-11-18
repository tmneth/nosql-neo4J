2.1. Surasti esybes pagal savybę (pvz. rasti asmenį pagal asmens kodą, rasti banko sąskaitą pagal numerį).
2.2. Surasti esybes pagal ryšį (pvz. banko sąskaitas priklausančias asmeniui, banko korteles susietas su konkretaus asmens sąskaitomis).
2.3. Surasti esybes susietas giliais sąryšiais (pvz. draugų draugus, visus kelius tarp Vilniaus ir Klaipėdos; visus autobusus kuriais galima nuvažiuoti iš stotelės X į stotelę Y).
2.4. Surasti trumpiausią kelią įvertinant svorius (pvz. surasti trumpiausią kelią tarp Vilniaus ir Klaipėdos; surasti pigiausią būdą konvertuoti iš valiutos X į valiutą Y, kuomet turima visų bankų konversijos informacija ir optimalus būdas, gali būti atlikti kelis žingsnius).
2.5. Agreguojami duomenys (pvz. kaip 2.4, tik surasti kelio ilgį ar konversijos kainą). Nenaudokite trumpiausio kelio.

- [ ] Find provider by ID

```
MATCH (c:ServiceProvider {provider_id: 'SP001'})
RETURN c;
```

- [ ] Find all cell towers of a service provider

```
MATCH (sp:ServiceProvider {provider_id: 'SP001'})-[:USES]->(ct:CellTower)
RETURN ct;

```

- [ ] Find all the data centers connected to the cell towers

```
MATCH (sp:ServiceProvider {provider_id: 'SP001'})-[:USES]->(ct:CellTower)-[:CONNECTS_TO]->(dc:DataCenter)
RETURN dc;

```
