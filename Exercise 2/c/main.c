// Compile with `gcc main.c -Wall -std=gnu99 -lpthread`, or use the makefile or gcc main.c -o main -lsocket
// The executable will be named `main` if you use the makefile, or `a.out` if you use gcc directly

#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <io.h>
#include <winsock2.h>


#define PORT 20000
#define SERVER_IP "10.22.77.124"

int main() {
    WSADATA wsaData;

    if (WSAStartup(MAKEWORD(2, 2), &wsaData) != 0) {
        perror("WSAStartup failed");
        return EXIT_FAILURE;
    }

    int client_socket;
    struct sockaddr_in server_address;

    // Create socket
    if ((client_socket = socket(AF_INET, SOCK_DGRAM, 0)) == -1) {
        perror("Socket creation failed");
        exit(EXIT_FAILURE);
    }

    // Set up server address structure
    server_address.sin_family = AF_INET;
    server_address.sin_port = htons(PORT);
    server_address.sin_addr.s_addr = inet_addr(SERVER_IP);

    // Perform data exchange with the server - You can add your own code here
    const char *message = "Hello, UDP Server!";
    if (sendto(client_socket, message, strlen(message), 0, (struct sockaddr *)&server_address, sizeof(server_address)) == -1) {
        perror("Send failed");
        exit(EXIT_FAILURE);
    }

    printf("Data sent to server %s:%d\n", SERVER_IP, PORT);

    close(client_socket); // Close the client socket
    return 0;
}




