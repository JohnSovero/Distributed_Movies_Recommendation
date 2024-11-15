import { Injectable } from '@angular/core';
import { HttpClient, HttpHeaders } from '@angular/common/http';
import { Observable, Subject, forkJoin } from 'rxjs';
import { Movie } from '../models/movie.model';
import { UserService } from './user.service';
import { map, switchMap } from 'rxjs/operators';
import { environment } from '../../../environments/environment';

@Injectable({
  providedIn: 'root'
})
export class WebsocketService {
  private ws!: WebSocket;
  private messagesSubject = new Subject<Movie[]>();
  public messages$ = this.messagesSubject.asObservable();
  private userId: number;

  private readonly TMDB_BASE_URL = 'https://api.themoviedb.org/3/movie';
  private readonly IMAGE_BASE_URL = 'https://image.tmdb.org/t/p/original';
  private readonly API_OPTIONS = {
    headers: new HttpHeaders({
      accept: 'application/json',
      Authorization: `Bearer ${environment.tmdbApiReadToken}`  // Replace with your actual API key
    })
  };

  constructor(private userService: UserService, private http: HttpClient) {
    this.userId = this.userService.getUserId();
  }

  connect(url: string) {
    const fullUrl = `${url}?userId=${this.userId}`;
    this.ws = new WebSocket(fullUrl);

    this.ws.onopen = () => {
      console.log('Connected to WebSocket.');
    };

    this.ws.onmessage = (event) => {
      const data = JSON.parse(event.data) as Movie[];

      // Fetch and assign poster URLs
      const movieObservables = data.map((movie: Movie) => this.fetchMoviePoster(movie));
      
      // Wait for all poster URLs to be retrieved, then emit the updated movies array
      forkJoin(movieObservables).subscribe(updatedMovies => {
        this.messagesSubject.next(updatedMovies);
      });
    };

    this.ws.onerror = (error) => {
      console.error('WebSocket error:', error);
    };

    this.ws.onclose = () => {
      console.log('WebSocket connection closed.');
    };
  }

  private fetchMoviePoster(movie: Movie): Observable<Movie> {
    const url = `${this.TMDB_BASE_URL}/${movie.tmdb_link}/images`;
    return this.http.get<any>(url, this.API_OPTIONS).pipe(
      map(response => {
        const backdropPath = response.backdrops?.[0]?.file_path;
        movie.poster = backdropPath ? `${this.IMAGE_BASE_URL}${backdropPath}` : '';
        console.log(movie);
        return movie;
      })
    );
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