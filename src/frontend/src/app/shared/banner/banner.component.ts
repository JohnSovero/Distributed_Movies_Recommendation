import { Component } from '@angular/core';
import { Subscription } from 'rxjs';
import { environment } from '../../../environments/environment';
import { WebsocketService } from '../../core/services/websocket.service';
import { HttpClient } from '@angular/common/http';
import { Movie } from '../../core/models/movie.model';
import { CommonModule } from '@angular/common';

@Component({
  selector: 'app-banner',
  standalone: true,
  imports: [CommonModule],
  templateUrl: './banner.component.html',
  styleUrls: ['./banner.component.css']
})
export class BannerComponent {
  bannerRecommendations: Movie[] = [];
  private bannerSubscription!: Subscription;
  private apiKey = environment.tmdbApiKey;
  currentMovieIndex: number = 0;
  showRecommendations: boolean = false;

  constructor(
    private websocketService: WebsocketService,
    private httpClient: HttpClient
  ) {}

  ngOnInit(): void {
    this.connect();
  }

  connect() {
    this.websocketService.connect('ws://localhost:9015/recommendations/above-average');

    this.bannerSubscription = this.websocketService.messages$.subscribe(
      (data) => {
        this.showRecommendations = false;
        this.bannerRecommendations = data;
        this.updatePosterPaths();
        console.log(this.bannerRecommendations);
        this.startBannerRotation();
        this.showRecommendations = true;
      }
    );
  }

  startBannerRotation() {
    setInterval(() => {
      this.currentMovieIndex = (this.currentMovieIndex + 1) % this.bannerRecommendations.length;
    }, 15000); // 5 seconds interval
  }

  updatePosterPaths() {
    this.bannerRecommendations.forEach((movie) => {
      const tmdbId = movie.tmdb_link;
      this.fetchMoviePosterPath(tmdbId).subscribe(
        (movieDetails) => {
          if (movieDetails.poster_path == null) {
            movie.poster = 'https://www.serieslike.com/img/shop_01.png';
          } else {
            movie.poster = `https://image.tmdb.org/t/p/w500${movieDetails.poster_path}`;
          }
          console.log('Poster URL:', movie.poster);
        },
        (error) => {
          movie.poster = 'https://www.serieslike.com/img/shop_01.png';
          console.log('Error fetching poster, using default:', movie.poster);
        }
      );
    });
  }

  fetchMoviePosterPath(tmdbId: string) {
    const url = `https://api.themoviedb.org/3/movie/${tmdbId}?api_key=${this.apiKey}`;
    return this.httpClient.get<any>(url);
  }

  ngOnDestroy() {
    this.bannerSubscription.unsubscribe();
    this.websocketService.disconnect();
  }
}