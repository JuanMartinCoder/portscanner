# NetPort

NetPort is a network scanning tool built from scratch in Go to analyze hosts, identify open ports, and determine their associated services and protocols. 
This project was born out of an interest in understanding the low-level workings of network communications and the foundations of cybersecurity.

### What does NetPort do? 🧐

This scanner allows you to get a quick overview of a host's "attack surface" by answering questions like:

* Which ports are open and accepting connections?
* What well-known service (HTTP, SSH, FTP, etc.) is expected to be running on that port?

### Implemented Features ✨

* **Open Port Detection:** Uses TCP Connections to check the status of ports within a specified range.
* **Protocol/Service Mapping:** Integrates an internal csv to correlate each open port with its most likely service.
* **Clean Output:** Displays the results in a clear and organized format in the terminal.
* **MultiThreading:** Uses Go routines for scanning faster.

### Future Improvements (Roadmap) 🗺️

* [ ] Add banner grabbing for more accurate service version identification.
* [ ] Support for different scan types (e.g., SYN, FIN).


<img width="892" height="241" alt="image" src="https://github.com/user-attachments/assets/bb6f06dd-b9b1-4a7c-b282-cfdbacdfd9f7" />
