
# **Fleetfy - Fleet Management and Logistics Platform**

Fleetfy is a web-based logistics and fleet management application. It connects users who need transportation services with drivers, while allowing administrators to manage the fleet of vehicles, monitor driver activity, and oversee bookings. The app provides real-time tracking, vehicle assignment, and role-based access control for users, drivers, and admins.
## **Important Links**
-   [Application Documentation](https://www.notion.so/Fleetfy-Documentation-1239888512ac81ccb81de7ac839dd3a0)
-   [Server Design and Performance Documentation](https://vaulted-hamster-fc2.notion.site/Fleetfy-Server-Implementation-and-Performance-Documentation-1239888512ac8003b8c8cb8b76816f9b)
-   [Database Design Document](https://www.notion.so/Database-Design-1239888512ac80bbb82beba69c466861)
-   [ER Digram](https://www.notion.so/ER-Diagram-of-our-Database-1239888512ac8093bc07dc0378e79f62)
-   [Future Scope for Fleetfy](https://vaulted-hamster-fc2.notion.site/Future-Scope-for-Fleetfy-1239888512ac805680dcfb16f5f73975)
-   [Application Demo Video](https://drive.google.com/file/d/1fsz9QHP2fIVyEmc7CuBTjipj1ivQxYee/view?usp=sharing)

## **Table of Contents**


-   [Tech Stack](#tech-stack)
-   [Installation](#installation)
-   [Running the Application](#running-the-application)

## **Tech Stack**

### **Backend**:

-   **Go (Golang)**: High-performance, statically typed backend.
-   **PostgreSQL**: Relational database for storing users, bookings, vehicles, and drivers.
-   **Gorilla Mux**: Router used for managing API endpoints.


### **Frontend**:

-   **React**: For building the user interface.
-   **Tailwind CSS**: For styling and responsive design.
-   **Axios**: For making HTTP requests from the frontend.

### **Other**:

-   **Nginx**: Used for load balancing multiple instances of the backend.
-   **Postman**: For API testing.



## **Installation**

### Prerequisites:

-   **Go** (v1.16 or later)
-   **PostgreSQL** (v12 or later)
-   **Node.js** (v14 or later for React)

### Clone the Repository:



`
git clone https://github.com/ravjot07/fleetfy.git
`





## **Running the Application**

### Backend (Go):

1.  Install Go dependencies:
    
     
    `go mod download` 
    
2.  Start the backend server:
    
    
    `go run main.go` 
    

### Frontend (React):

1.  Navigate to the frontend directory and install dependencies:
    

    `cd frontend
    npm install` 
    
2.  Start the development server:
    
    
    `npm run dev` 
    

The backend will run at `http://localhost:8080`, and the frontend at `http://localhost:5173`.

