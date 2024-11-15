import { HttpClient } from '@angular/common/http';
import { Component, ElementRef, HostListener, ViewChild } from '@angular/core';
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
  @ViewChild('genreRec') genreRec: ElementRef | undefined;
  private isMouseDown: boolean = false;
  private startX: number = 0;
  private scrollLeft: number = 0;

  constructor(
    private dataService: DataService,
    private httpClient: HttpClient
  ) {}

  ngOnInit(): void {
    this.onButtonSelect(this.selectedGenre);
    this.getRecommendations();
  }

  @HostListener('mousedown', ['$event'])
  onMouseDown(event: MouseEvent): void {
    if (this.genreRec?.nativeElement) {
      this.isMouseDown = true;
      this.startX = event.pageX - this.genreRec.nativeElement.offsetLeft;
      this.scrollLeft = this.genreRec.nativeElement.scrollLeft;
    }
  }

  @HostListener('mouseup')
  onMouseUp(): void {
    this.isMouseDown = false;
  }

  @HostListener('mousemove', ['$event'])
  onMouseMove(event: MouseEvent): void {
    if (!this.isMouseDown || !this.genreRec?.nativeElement) return;
    const x = event.pageX - this.genreRec.nativeElement.offsetLeft;
    const walk = (x - this.startX) * 2; // Adjust the 2 for scroll speed
    this.genreRec.nativeElement.scrollLeft = this.scrollLeft - walk;
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