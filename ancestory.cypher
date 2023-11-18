MATCH (n) DETACH DELETE n;

CREATE (r1:Region {region_id: 'R001', name: 'Baltic Region'}),
       (r2:Region {region_id: 'R002', name: 'Nordic Region'}),
       (r3:Region {region_id: 'R003', name: 'Mediterranean Region'}),
       (r4:Region {region_id: 'R004', name: 'Central European Region'}),
       (r5:Region {region_id: 'R005', name: 'Eastern European Region'}),
       (r6:Region {region_id: 'R006', name: 'Western European Region'}),
       (r7:Region {region_id: 'R007', name: 'South American Region'}),
    
       (p1:Person {person_id: 'P001', name: 'John'}),
       (p2:Person {person_id: 'P002', name: 'Alice'}),
       (p3:Person {person_id: 'P003', name: 'Bob'}),
       (p4:Person {person_id: 'P004', name: 'Charlie'}),
       (p5:Person {person_id: 'P005', name: 'Diane'}),
       (p6:Person {person_id: 'P006', name: 'Eva'}),

       (s1:Surname {surname_id: 'S001', surname: 'Smith'}),
       (s2:Surname {surname_id: 'S002', surname: 'Johnson'}),
       (s3:Surname {surname_id: 'S003', surname: 'Williams'}),

       (p1)-[:HAS_SURNAME]->(s1),
       (p2)-[:HAS_SURNAME]->(s1),
       (p3)-[:HAS_SURNAME]->(s2),
       (p4)-[:HAS_SURNAME]->(s3),
       (p5)-[:HAS_SURNAME]->(s3),
       (p6)-[:HAS_SURNAME]->(s1),

       (p1)-[:FROM_REGION {probability: 0.4}]->(r1),
       (p1)-[:FROM_REGION {probability: 0.3}]->(r2),
       (p1)-[:FROM_REGION {probability: 0.2}]->(r3),
       (p1)-[:FROM_REGION {probability: 0.1}]->(r4),

       (p2)-[:FROM_REGION {probability: 0.25}]->(r2),
       (p2)-[:FROM_REGION {probability: 0.25}]->(r4),
       (p2)-[:FROM_REGION {probability: 0.25}]->(r5),
       (p2)-[:FROM_REGION {probability: 0.25}]->(r6),

       (p3)-[:FROM_REGION {probability: 0.4}]->(r3),
       (p3)-[:FROM_REGION {probability: 0.3}]->(r1),
       (p3)-[:FROM_REGION {probability: 0.2}]->(r5),
       (p3)-[:FROM_REGION {probability: 0.1}]->(r6),

       (p4)-[:FROM_REGION {probability: 0.35}]->(r4),
       (p4)-[:FROM_REGION {probability: 0.35}]->(r2),
       (p4)-[:FROM_REGION {probability: 0.2}]->(r6),
       (p4)-[:FROM_REGION {probability: 0.1}]->(r7),

       (p5)-[:FROM_REGION {probability: 0.4}]->(r1),
       (p5)-[:FROM_REGION {probability: 0.3}]->(r4),
       (p5)-[:FROM_REGION {probability: 0.2}]->(r5),
       (p5)-[:FROM_REGION {probability: 0.1}]->(r6),

       (p6)-[:FROM_REGION {probability: 0.45}]->(r1),
       (p6)-[:FROM_REGION {probability: 0.35}]->(r3),
       (p6)-[:FROM_REGION {probability: 0.1}]->(r6),
       (p6)-[:FROM_REGION {probability: 0.1}]->(r7),

       (p1)-[:RELATED_TO {probability: 0.05}]->(p2),
       (p2)-[:RELATED_TO {probability: 0.03}]->(p3),
       (p3)-[:RELATED_TO {probability: 0.02}]->(p1),
       (p4)-[:RELATED_TO {probability: 0.04}]->(p5),
       (p5)-[:RELATED_TO {probability: 0.01}]->(p6),
       (p6)-[:RELATED_TO {probability: 0.03}]->(p4);