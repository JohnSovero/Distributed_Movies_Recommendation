import { Injectable } from '@angular/core';
import { Observable, Subject } from 'rxjs';
import { Movie } from '../models/movie.model';
import { UserService } from './user.service';  // Import the UserService

@Injectable({
  providedIn: 'root'
})
export class WebsocketService {
  private ws!: WebSocket;
  private messagesSubject = new Subject<Movie[]>();
  public messages$ = this.messagesSubject.asObservable();
  private userId: number;

  constructor(private userService: UserService) {
    // Access the global userId from UserService
    this.userId = this.userService.getUserId();
  }

  connect(url: string) {
    const fullUrl = `${url}?userId=${this.userId}`; // Append the userId to the WebSocket URL
    this.ws = new WebSocket(fullUrl);

    this.ws.onopen = () => {
      console.log('Connected to WebSocket.');
    };

    this.ws.onmessage = (event) => {
      const data = JSON.parse(event.data);

      // Set poster to null initially
      data.forEach((movie: Movie) => {
        movie.poster = ''; // Set poster to null initially
      });

      this.messagesSubject.next(data);
    };

    this.ws.onerror = (error) => {
      console.error('WebSocket error:', error);
    };

    this.ws.onclose = () => {
      console.log('WebSocket connection closed.');
    };
  }

  sendMessage(message: any) {
    if (this.ws.readyState === WebSocket.OPEN) {
      this.ws.send(JSON.stringify(message));
    } else {
      console.error('WebSocket connection is not open.');
    }
  }

  disconnect() {
    if (this.ws) {
      this.ws.close();
    }
  }
}