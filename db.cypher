// Create Service Providers
CREATE (sp1:ServiceProvider {provider_id: 'SP001', name: 'TeleFast', service_area: 'North Region'})
CREATE (sp2:ServiceProvider {provider_id: 'SP002', name: 'QuickCall', service_area: 'South Region'})

// Create Cell Towers
CREATE (ct1:CellTower {tower_id: 'CT001', location: 'LocationA', capacity: 100, operational_status: 'Active'})
CREATE (ct2:CellTower {tower_id: 'CT002', location: 'LocationB', capacity: 150, operational_status: 'Active'})
CREATE (ct3:CellTower {tower_id: 'CT003', location: 'LocationC', capacity: 80, operational_status: 'Maintenance'})

// Create Data Centers
CREATE (dc1:DataCenter {center_id: 'DC001', location: 'Central Hub', max_load: 1000, power_backup_status: 'Yes'})
CREATE (dc2:DataCenter {center_id: 'DC002', location: 'Backup Site', max_load: 500, power_backup_status: 'No'})

// Create Customers
CREATE (c1:Customer {customer_id: 'C001', name: 'John Doe', address: '123 Main St'})
CREATE (c2:Customer {customer_id: 'C002', name: 'Jane Smith', address: '456 Broad Ave'})
CREATE (c3:Customer {customer_id: 'C003', name: 'Alice Johnson', address: '789 Oak Ln'})

// Create Relationships between Service Providers and Cell Towers
CREATE (sp1)-[:USES]->(ct1)
CREATE (sp1)-[:USES]->(ct2)
CREATE (sp2)-[:USES]->(ct3)

// Create Relationships between Cell Towers and Data Centers
CREATE (ct1)-[:CONNECTS_TO {bandwidth: '10Gbps'}]->(dc1)
CREATE (ct2)-[:CONNECTS_TO {bandwidth: '1Gbps'}]->(dc1)
CREATE (ct3)-[:CONNECTS_TO {bandwidth: '100Mbps'}]->(dc2)

// Create Relationships between Customers and Service Providerst
CREATE (c1)-[:SUBSCRIBED_TO]->(sp1)
CREATE (c2)-[:SUBSCRIBED_TO]->(sp1)
CREATE (c3)-[:SUBSCRIBED_TO]->(sp2)
