import { Component } from '@angular/core';
import { Subscription } from 'rxjs';
import { environment } from '../../../environments/environment';
import { WebsocketService } from '../../core/services/websocket.service';
import { HttpClient } from '@angular/common/http';
import { Movie } from '../../core/models/movie.model';
import { CommonModule } from '@angular/common';
import { LoaderComponent } from '../loader/loader.component';

@Component({
  selector: 'app-banner',
  standalone: true,
  imports: [CommonModule, LoaderComponent],
  templateUrl: './banner.component.html',
  styleUrls: ['./banner.component.css']
})
export class BannerComponent {
  // bannerRecommendations: Movie[] = [];
  bannerRecommendations: Movie[] = [
    {
      id: 1,
      title: 'Toy Story',
      genres: ["Adventure", "Animation", "Children", "Comedy", "Fantasy"],
      imdb_link: "114709",
      tmdb_link: "862",
      year: 1995,
      overview: "Led by Woody, Andy's toys live happily in his room until Andy's birthday brings Buzz Lightyear onto the scene. Afraid of losing his place in Andy's heart, Woody plots against Buzz. But when circumstances separate Buzz and Woody from their owner, the duo eventually learns to put aside their differences.",
      vote_avg: "7.7",
      poster: 'https://image.tmdb.org/t/p/original/lxD5ak7BOoinRNehOCA85CQ8ubr.jpg'
    },
    {
      id: 2,
      title: 'Jumanji',
      genres: ["Adventure", "Children", "Fantasy"],
      imdb_link: "113497",
      tmdb_link: "8844",
      year: 1995,
      overview: "When siblings Judy and Peter discover an enchanted board game that opens the door to a magical world, they unwittingly invite Alan -- an adult who's been trapped inside the game for 26 years -- into their living room. Alan's only hope for freedom is to finish the game, which proves risky as all three find themselves running from giant rhinoceroses, evil monkeys and other terrifying creatures.",
      vote_avg: "6.9",
      poster: 'https://image.tmdb.org/t/p/original/okURFlZWBlaq88NeyPuOqXxX7g.jpg'
    }
  ];
  private bannerSubscription!: Subscription;
  selectedIndex: number = 0;
  showRecommendations: boolean = true;
  private rotationTimeout!: ReturnType<typeof setTimeout>;

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
        // this.showRecommendations = false;
        this.bannerRecommendations = data;
        console.log(this.bannerRecommendations);
        this.startBannerRotation();
        this.showRecommendations = true;
      }
    );
  }

  startBannerRotation() {
    const rotateBanner = () => {
      if (this.bannerRecommendations.length > 0) {
        this.selectedIndex = (this.bannerRecommendations.length + this.selectedIndex + 1) % this.bannerRecommendations.length;
        console.log(this.selectedIndex);
      }

      // Schedule the next rotation with a dynamically adjusted interval
      const nextInterval = 15000;
      this.rotationTimeout = setTimeout(rotateBanner, nextInterval);
    };

    rotateBanner();
  }

  ngOnDestroy() {
    clearTimeout(this.rotationTimeout);
    this.bannerSubscription.unsubscribe();
    this.websocketService.disconnect();
  }
}