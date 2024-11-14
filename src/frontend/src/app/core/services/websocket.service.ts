import { Injectable } from '@angular/core';
import { Observable, Subject } from 'rxjs';
import { Movie } from '../models/movie.model';

@Injectable({
  providedIn: 'root'
})
export class WebsocketService {
  private ws!: WebSocket;
  private messagesSubject = new Subject<Movie[]>();
  public messages$ = this.messagesSubject.asObservable();

  constructor() { }

  connect(url: string) {
    this.ws = new WebSocket(url);

    this.ws.onopen = () => {
      console.log('Connected to WebSocket.');
    };

    this.ws.onmessage = (event) => {
      const data = JSON.parse(event.data);
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