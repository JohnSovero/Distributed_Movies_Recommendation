import { HttpClient } from '@angular/common/http';
import { Component, ElementRef, HostListener, ViewChild } from '@angular/core';
import { DataService } from '../../core/services/data.service';
import { Movie } from '../../core/models/movie.model';
import { environment } from '../../../environments/environment';
import { MatChipsModule } from '@angular/material/chips';
import { CommonModule } from '@angular/common';
import { MatButtonModule } from '@angular/material/button';
import { LoaderComponent } from '../loader/loader.component';
import { MatBottomSheet } from '@angular/material/bottom-sheet';
import { BottomSheetComponent } from '../bottom-sheet/bottom-sheet.component';
import { MatSliderModule } from '@angular/material/slider';
import { FormsModule } from '@angular/forms';

@Component({
  selector: 'app-genre-rec',
  standalone: true,
  imports: [MatChipsModule, CommonModule, MatButtonModule, LoaderComponent, MatSliderModule, FormsModule],
  templateUrl: './genre-rec.component.html',
  styleUrls: ['./genre-rec.component.css']
})
export class GenreRecComponent {
  genres: string[] = ['All', 'Action', 'Animation', 'Comedy', 'Drama', 'Horror', 'Musical', 'Children', 'Romance'];
  selectedGenre: string = 'All';
  selectedGenreTemp: string = 'All';
  genreRecommendations: Movie[] = [];
  showRecommendations: boolean = false;
  private apiKey = environment.tmdbApiKey;
  sliderValue: number = 50;

  constructor(
    private dataService: DataService,
    private httpClient: HttpClient,
    private _bottomSheet: MatBottomSheet
  ) {}

  ngOnInit(): void {
    this.getRecommendations(this.sliderValue);
  }

  openBottomSheet(movie: Movie): void {
    console.log('Opening bottom sheet with movie:', movie);
    const bottomSheetRef = this._bottomSheet.open(BottomSheetComponent, {data: movie});
  }

  getRecommendations(numRec: number): void {
    this.showRecommendations = false;
    this.dataService.getRecommendations(this.selectedGenreTemp, numRec).subscribe((movies) => {
      this.genreRecommendations = movies;
      console.log('Movies arrived! time:', new Date().toLocaleTimeString());
      this.updatePosterPaths();
      console.log('Posters updated! time:', new Date().toLocaleTimeString());
      this.showRecommendations = true;
    });
  }

  onButtonSelect(genre: string): void {
    if (this.selectedGenreTemp !== genre) {
      this.selectedGenreTemp = genre;
    }
  }

  recommendMovies(): void {
    this.selectedGenre = this.selectedGenreTemp;
    this.getRecommendations(this.sliderValue);
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

  formatLabel(value: number): string {
    if (value >= 1000) {
      return Math.round(value / 1000) + 'k';
    }
    return value.toString();
  }

}