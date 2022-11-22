- [Data Modeling](#data-modeling)
  - [Data Structures](#data-structures)
  - [Directed Propeerty Graph](#directed-propeerty-graph)
- [VID](#vid)
  - [Features](#features)
  - [VID Operation](#vid-operation)

# Data Modeling

## Data Structures

> NebulaGraph data model uses six data structures to store data. They are graph spaces, vertices, edges, tags, edge types and properties.

- **Graph spaces**: Graph spaces are used to isolate data from different teams or programs. Data stored in different graph spaces are securely isolated. Storage replications, privileges, and partitions can be assigned.
  
- **Vertices**: Vertices are used to store entities.

    In NebulaGraph, vertices are identified with vertex identifiers (i.e. `VID`). The `VID` must be unique in the same graph space. `VID` should be `int64`, or `fixed_string(N)`.

    ~~~warning
    In NebulaGraph 2.x a vertex must have at least one tag. And in NebulaGraph 3.3.0, a tag is not required for a vertex.
    ~~~

- **Edges**: Edges are used to connect vertices. An edge is a connection or behavior between two vertices.
  
  - There can be multiple edges between two vertices.
  - Edges are directed. -> identifies the directions of edges. Edges can be traversed in either direction.
  - An edge is identified uniquely with `<a source vertex, an edge type, a rank value, and a destination vertex>`. **Edges have no EID**.
  - An edge must have one and only one edge type.
  - The rank value is an immutable user-assigned 64-bit signed integer. It identifies the edges with the same edge type between two vertices. Edges are sorted by their rank values. The edge with the greatest rank value is listed first. The default rank value is zero.

- **Tags**: Tags are used to categorize vertices. Vertices that have the same tag share the same definition of properties.

- **Edge types**: Edge types are used to categorize edges. Edges that have the same edge type share the same definition of properties.

- **Properties**: Properties are key-value pairs. Both vertices and edges are containers for properties.

```note
Tags and Edge types are similar to "vertex tables" and "edge tables" in the relational databases.
```

## Directed Propeerty Graph

NebulaGraph stores data in directed property graphs. A directed property graph has a set of vertices connected by directed edges. Both vertices and edges can have properties. A directed property graph is represented as:

`G = < V, E, PV, PE >`

  - V is a set of vertices.
  - E is a set of directed edges.
  - PV is the property of vertices.
  - PE is the property of edges.

The following table is an example of the structure of the basketball player dataset. We have two types of vertices, that is player and team, and two types of edges, that is serve and follow.

| Element   | Name   | Property name (Data type)       | Description                                                                                                                                              |
| --------- | ------ | ------------------------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Tag       | player | name (string), age (int)        | Represents players in the team.                                                                                                                          |
| Tag       | team   | name (string)                   | Represents the teams.                                                                                                                                    |
| Edge type | serve  | start_year (int),end_year (int) | Represents actions taken by players in the team. An action links a player with a team, and the direction is from a player to a team.                     |
| Edge type | follow | degree (int)                    | Represents actions taken by players in the team. An action links a player with another player, and the direction is from one player to the other player. |

```info
NebulaGraph 3.3.0 allows dangling edges. Therefore, when adding or deleting, you need to ensure the corresponding source vertex and destination vertex of an edge exist. For details, see INSERT VERTEX, DELETE VERTEX, INSERT EDGE, and DELETE EDGE.

The MERGE statement in openCypher is not supported.
```

# VID

```info
In NebulaGraph, a vertex is uniquely identified by its ID, which is called a VID or a Vertex ID.
```ne

## Features

- The data types of VIDs are restricted to `FIXED_STRING`(<N>) or `INT64`. One graph space can only select one VID type.
- A VID in a graph space is `unique`. It functions just as a `primary key` in a relational database. VIDs in different graph spaces are independent.
- **The VID generation method must be set by users, because NebulaGraph does not provide auto increasing ID, or UUID.**
- Vertices with the same VID will be identified as the same one. For example:

  - A VID is the unique identifier of an entity, like a person's ID card number. A tag means the type of an entity, such as driver, and boss. Different tags define two groups of different properties, such as driving license number, driving age, order amount, order taking alt, and job number, payroll, debt ceiling, business phone number.
  - When two `INSERT` statements (neither uses a parameter of IF NOT EXISTS) with the same VID and tag are operated at the same time, **the latter INSERT will overwrite the former**.
  - When two `INSERT` statements with the same VID but different tags, like **TAG A** and **TAG B**, are operated at the same time, the operation of **Tag A will not affect Tag B**.

- VIDs will usually be indexed and stored into memory (in the way of `LSM-tree`). Thus, direct access to VIDs enjoys peak performance.

## VID Operation

NebulaGraph 1.x only supports INT64 while NebulaGraph 2.x supports INT64 and FIXED_STRING(<N>). In CREATE SPACE, VID types can be set via vid_type.
id() function can be used to specify or locate a VID.
LOOKUP or MATCH statements can be used to find a VID via property index.
Direct access to vertices statements via VIDs enjoys peak performance, such as DELETE xxx WHERE id(xxx) == "player100" or GO FROM "player100". Finding VIDs via properties and then operating the graph will cause poor performance, such as LOOKUP | GO FROM $-.ids, which will run both LOOKUP and | one more time.