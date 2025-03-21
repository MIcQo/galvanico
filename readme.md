# Galvanico

[![Contributions Welcome](https://img.shields.io/badge/Contributions-Welcome-brightgreen.svg)](CONTRIBUTING.md)

Galvanico is an open-source, browser-based strategy game inspired by Ikariam, but set in the Industrial Age. Players develop their own industrial cities, research new technologies, manage resources, and engage in trade and diplomacy in a world powered by electricity and innovation.

## Features

* **Industrial Age Setting:** Immerse yourself in a world of steam engines, early electricity, and burgeoning industries.
* **City Building:** Construct and manage various industrial buildings, including factories, power plants, and research labs.
* **Research & Technology:** Unlock new technologies to improve your city's efficiency, military strength, and economic power.
* **Resource Management:** Gather and manage resources like coal, iron, and electricity to fuel your industrial empire.
* **Trade & Diplomacy:** Interact with other players through trade agreements, alliances, and diplomacy.
* **Military Expansion:** Build and command industrial-era military units to defend your city and conquer new territories.
* **Electricity System:** A core mechanic, managing and generating electricity to power your industrial buildings.
* **Open Source:** Contribute to the development and shape the future of Galvanico.

## Getting Started

### Prerequisites

* Go (latest stable version)
* Node.js (latest LTS recommended)
* npm (or yarn)
* A modern web browser

### Installation

1.  Clone the repository:

    ```bash
    git clone [https://github.com/yourusername/galvanico.git](https://www.google.com/search?q=https://github.com/yourusername/galvanico.git)
    cd galvanico
    ```

2.  Install backend dependencies (Go):

    ```bash
    go mod tidy
    ```

3.  Install frontend dependencies (Vue 3):

    ```bash
    cd client
    npm install # or yarn install
    cd ..
    ```

4.  Start the backend server:

    ```bash
    go run server/main.go
    ```

5.  Start the frontend development server:

    ```bash
    cd client
    npm run dev # or yarn dev
    ```

6.  Open your browser and navigate to `http://localhost:5173`. (or the port that Vue dev server outputs)

## Development

We welcome contributions from the community! Please read our [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines on how to contribute.

### Project Structure
galvanico/
├── client/          # Frontend code (Vue 3, etc.)
├── server/          # Backend code (Go)
├── common/          # Shared code between client and server (if any)
├── assets/          # Images, sounds, and other assets
├── docs/            # Documentation
├── CONTRIBUTING.md  # Contribution guidelines
├── LICENSE          # License information
└── README.md        # This file

### Technologies Used

* **Frontend:** [Vue 3](https://vuejs.org/) , [TypeScript](https://www.typescriptlang.org/)
* **Backend:** [Go](https://golang.org/)
* **Database:** [PostgreSQL](https://www.postgresql.org/) (or your preferred database)

## Contributing

We encourage you to contribute to Galvanico! Here's how you can get involved:

* **Report Bugs:** If you find a bug, please open an issue on GitHub.
* **Suggest Features:** Share your ideas for new features and improvements.
* **Submit Pull Requests:** Contribute code, documentation, or other improvements.
* **Help with Design:** Contribute to the game's visual design and user interface.
* **Translate the Game:** Help us make Galvanico accessible to a wider audience.

Please read our [CONTRIBUTING.md](CONTRIBUTING.md) for detailed guidelines.

## License

Galvanico is released under the [MIT License](LICENSE).

## Acknowledgments

* Inspired by [Ikariam](https://www.ikariam.com/).
* Thanks to all contributors and the open-source community.

## Contact

For any questions or inquiries, please open an issue on GitHub.

---

**Let's build the industrial age together!**