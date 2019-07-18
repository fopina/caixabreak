from django.conf.urls import patterns, url
from django.views.generic import TemplateView
import views

urlpatterns = patterns('',
    url(r'^$', views.index, name='index'),
    url(r'^info/$', views.info, name='info'),
    url(r'^appcache/', TemplateView.as_view(template_name="cbreak/appcache.html", content_type='text/cache-manifest'), name='appcache'),
)
