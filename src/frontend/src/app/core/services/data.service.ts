import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { UserService } from './user.service';
import { map, Observable } from 'rxjs';
import { Movie } from '../models/movie.model';

@Injectable({
  providedIn: 'root',
})
export class DataService {
  private userId: number;
  private baseUrl = 'http://localhost:9015';

  constructor(
    private http: HttpClient,
    private userService: UserService,
  ) {
    this.userId = this.userService.getUserId();
  }

  getRecommendations(genre: string, numRec: number): Observable<Movie[]> {
    return this.http.get<Movie[]>(`${this.baseUrl}/recommendations/${numRec}/genres/${genre}/users/${this.userId}`).pipe(
      map((movies) => movies.map(movie => ({ ...movie, poster: '' })))
    );
  }
}