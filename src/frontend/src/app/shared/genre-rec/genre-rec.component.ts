import { HttpClient } from '@angular/common/http';
import { Component } from '@angular/core';
import { DataService } from '../../core/services/data.service';
import { Movie } from '../../core/models/movie.model';
import { environment } from '../../../environments/environment';
import { MatChipsModule } from '@angular/material/chips';
import { CommonModule } from '@angular/common';
import { MatButtonModule } from '@angular/material/button';
import { LoaderComponent } from '../loader/loader.component';

@Component({
  selector: 'app-genre-rec',
  standalone: true,
  imports: [MatChipsModule, CommonModule, MatButtonModule, LoaderComponent],
  templateUrl: './genre-rec.component.html',
  styleUrls: ['./genre-rec.component.css']
})
export class GenreRecComponent {
  genres: string[] = ['Action', 'Animation', 'Comedy', 'Drama', 'Horror', 'Musical', 'Children', 'Romance', 'All'];
  selectedGenre: string = 'All';
  genreRecommendations: Movie[] = [];
  showRecommendations: boolean = false;
  private apiKey = environment.tmdbApiKey;

  constructor(
    private dataService: DataService,
    private httpClient: HttpClient
  ) {}

  ngOnInit(): void {
    this.onButtonSelect(this.selectedGenre);
    this.getRecommendations();
  }

  getRecommendations(): void {
    // console.log('Selected option:', this.selectedGenre);
    this.showRecommendations = false;
    this.dataService.getRecommendations(this.selectedGenre).subscribe((movies) => {
      this.genreRecommendations = movies;
      this.updatePosterPaths();
      this.showRecommendations = true;
    });
  }

  onButtonSelect(genre: string): void {
    if (this.selectedGenre !== genre && this.showRecommendations == true) {
      this.selectedGenre = genre;
      this.getRecommendations();
    }
  }

  updatePosterPaths() {
    this.genreRecommendations.forEach((movie) => {
      const tmdbId = movie.tmdb_link;
      this.fetchMoviePosterPath(tmdbId).subscribe(
        (movieDetails: any) => {
          if (movieDetails.poster_path == null) {
            movie.poster = 'https://www.serieslike.com/img/shop_01.png';
          } else {
            movie.poster = `https://image.tmdb.org/t/p/w500${movieDetails.poster_path}`;
          }
        },
        (error: any) => {
          movie.poster = 'https://www.serieslike.com/img/shop_01.png';
        }
      );
    });
  }

  fetchMoviePosterPath(tmdbId: string) {
    const url = `https://api.themoviedb.org/3/movie/${tmdbId}?api_key=${this.apiKey}`;
    return this.httpClient.get<any>(url);
  }
}