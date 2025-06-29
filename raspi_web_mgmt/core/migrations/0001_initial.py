# Generated by Django 5.2.3 on 2025-06-25 12:22

import django.db.models.deletion
from django.db import migrations, models


class Migration(migrations.Migration):

    initial = True

    dependencies = [
    ]

    operations = [
        migrations.CreateModel(
            name='AuthGroup',
            fields=[
                ('id', models.BigAutoField(auto_created=True, primary_key=True, serialize=False, verbose_name='ID')),
                ('name', models.CharField(max_length=150, unique=True)),
            ],
            options={
                'db_table': 'auth_group',
                'managed': False,
            },
        ),
        migrations.CreateModel(
            name='AuthGroupPermissions',
            fields=[
                ('id', models.BigAutoField(auto_created=True, primary_key=True, serialize=False, verbose_name='ID')),
            ],
            options={
                'db_table': 'auth_group_permissions',
                'managed': False,
            },
        ),
        migrations.CreateModel(
            name='AuthPermission',
            fields=[
                ('id', models.BigAutoField(auto_created=True, primary_key=True, serialize=False, verbose_name='ID')),
                ('codename', models.CharField(max_length=100)),
                ('name', models.CharField(max_length=255)),
            ],
            options={
                'db_table': 'auth_permission',
                'managed': False,
            },
        ),
        migrations.CreateModel(
            name='AuthUser',
            fields=[
                ('id', models.BigAutoField(auto_created=True, primary_key=True, serialize=False, verbose_name='ID')),
                ('password', models.CharField(max_length=128)),
                ('last_login', models.DateTimeField(blank=True, null=True)),
                ('is_superuser', models.BooleanField()),
                ('username', models.CharField(max_length=150, unique=True)),
                ('last_name', models.CharField(max_length=150)),
                ('email', models.CharField(max_length=254)),
                ('is_staff', models.BooleanField()),
                ('is_active', models.BooleanField()),
                ('date_joined', models.DateTimeField()),
                ('first_name', models.CharField(max_length=150)),
            ],
            options={
                'db_table': 'auth_user',
                'managed': False,
            },
        ),
        migrations.CreateModel(
            name='AuthUserGroups',
            fields=[
                ('id', models.BigAutoField(auto_created=True, primary_key=True, serialize=False, verbose_name='ID')),
            ],
            options={
                'db_table': 'auth_user_groups',
                'managed': False,
            },
        ),
        migrations.CreateModel(
            name='AuthUserUserPermissions',
            fields=[
                ('id', models.BigAutoField(auto_created=True, primary_key=True, serialize=False, verbose_name='ID')),
            ],
            options={
                'db_table': 'auth_user_user_permissions',
                'managed': False,
            },
        ),
        migrations.CreateModel(
            name='DjangoAdminLog',
            fields=[
                ('id', models.BigAutoField(auto_created=True, primary_key=True, serialize=False, verbose_name='ID')),
                ('object_id', models.TextField(blank=True, null=True)),
                ('object_repr', models.CharField(max_length=200)),
                ('action_flag', models.PositiveSmallIntegerField()),
                ('change_message', models.TextField()),
                ('action_time', models.DateTimeField()),
            ],
            options={
                'db_table': 'django_admin_log',
                'managed': False,
            },
        ),
        migrations.CreateModel(
            name='DjangoContentType',
            fields=[
                ('id', models.BigAutoField(auto_created=True, primary_key=True, serialize=False, verbose_name='ID')),
                ('app_label', models.CharField(max_length=100)),
                ('model', models.CharField(max_length=100)),
            ],
            options={
                'db_table': 'django_content_type',
                'managed': False,
            },
        ),
        migrations.CreateModel(
            name='DjangoMigrations',
            fields=[
                ('id', models.BigAutoField(auto_created=True, primary_key=True, serialize=False, verbose_name='ID')),
                ('app', models.CharField(max_length=255)),
                ('name', models.CharField(max_length=255)),
                ('applied', models.DateTimeField()),
            ],
            options={
                'db_table': 'django_migrations',
                'managed': False,
            },
        ),
        migrations.CreateModel(
            name='DjangoSession',
            fields=[
                ('session_key', models.CharField(max_length=40, primary_key=True, serialize=False)),
                ('session_data', models.TextField()),
                ('expire_date', models.DateTimeField()),
            ],
            options={
                'db_table': 'django_session',
                'managed': False,
            },
        ),
        migrations.CreateModel(
            name='Central',
            fields=[
                ('id', models.CharField(max_length=254, primary_key=True, serialize=False)),
                ('nome', models.CharField(max_length=254)),
                ('versao', models.CharField(max_length=254)),
                ('data_atualizacao', models.DateField(auto_now=True)),
            ],
            options={
                'db_table': 'central',
            },
        ),
        migrations.CreateModel(
            name='MicroControlador',
            fields=[
                ('id', models.CharField(max_length=254, primary_key=True, serialize=False)),
                ('tipo_controlador', models.CharField(max_length=254)),
                ('caminho_binario', models.CharField(max_length=254)),
                ('token', models.TextField(blank=True, null=True)),
                ('ip', models.CharField(max_length=254)),
                ('data_atualizacao', models.DateField(auto_now=True)),
            ],
            options={
                'db_table': 'micro_controlador',
            },
        ),
        migrations.CreateModel(
            name='HistoricoAtualizacao',
            fields=[
                ('id', models.BigAutoField(auto_created=True, primary_key=True, serialize=False, verbose_name='ID')),
                ('data_atualizacao', models.DateField(blank=True, null=True)),
                ('versao', models.CharField(max_length=254)),
                ('log', models.CharField(max_length=254)),
                ('micro', models.ForeignKey(on_delete=django.db.models.deletion.DO_NOTHING, to='core.microcontrolador')),
            ],
            options={
                'db_table': 'historico_atualizacao',
            },
        ),
        migrations.CreateModel(
            name='CentralMicro',
            fields=[
                ('id', models.BigAutoField(auto_created=True, primary_key=True, serialize=False, verbose_name='ID')),
                ('versao', models.TextField(blank=True, null=True)),
                ('data_atualizacao', models.DateField(auto_now=True)),
                ('central', models.ForeignKey(on_delete=django.db.models.deletion.DO_NOTHING, to='core.central')),
                ('micro', models.ForeignKey(on_delete=django.db.models.deletion.DO_NOTHING, to='core.microcontrolador')),
            ],
            options={
                'db_table': 'central_micro',
            },
        ),
        migrations.CreateModel(
            name='Topicos',
            fields=[
                ('id', models.BigAutoField(auto_created=True, primary_key=True, serialize=False, verbose_name='ID')),
                ('nome', models.TextField()),
                ('mensagem', models.TextField(blank=True, null=True)),
                ('data_atualizacao', models.DateField(auto_now=True)),
                ('central', models.ForeignKey(on_delete=django.db.models.deletion.DO_NOTHING, to='core.central')),
            ],
            options={
                'db_table': 'topicos',
            },
        ),
        migrations.CreateModel(
            name='Treino',
            fields=[
                ('id', models.BigAutoField(auto_created=True, primary_key=True, serialize=False, verbose_name='ID')),
                ('user_id', models.CharField(max_length=254)),
                ('plan_id', models.CharField(max_length=254)),
                ('main_muscles', models.CharField(max_length=254)),
                ('exercise_id', models.CharField(max_length=254)),
                ('gym_id', models.CharField(max_length=254)),
                ('summary', models.CharField(max_length=254)),
                ('workout_score', models.FloatField(blank=True, null=True)),
                ('exercise_score', models.FloatField(blank=True, null=True)),
                ('resting_time', models.FloatField(blank=True, null=True)),
                ('used_load', models.FloatField(blank=True, null=True)),
                ('failed_reps', models.IntegerField(blank=True, null=True)),
                ('total_reps', models.IntegerField(blank=True, null=True)),
                ('total_series', models.IntegerField(blank=True, null=True)),
                ('amplitude_score', models.FloatField(blank=True, null=True)),
                ('detailed_amplitude_score', models.FloatField(blank=True, null=True)),
                ('time_score', models.FloatField(blank=True, null=True)),
                ('detailed_time_score', models.FloatField(blank=True, null=True)),
                ('data_atualizacao', models.DateField(auto_now=True)),
                ('micro', models.ForeignKey(on_delete=django.db.models.deletion.DO_NOTHING, to='core.microcontrolador')),
            ],
            options={
                'db_table': 'treino',
            },
        ),
        migrations.CreateModel(
            name='VersaoOnline',
            fields=[
                ('id', models.BigAutoField(auto_created=True, primary_key=True, serialize=False, verbose_name='ID')),
                ('caminho_binario', models.CharField(max_length=254)),
                ('versao', models.CharField(max_length=254)),
                ('data_versao', models.DateField(blank=True, null=True)),
                ('tipo_controlador', models.CharField(max_length=254)),
                ('data_atualizacao', models.DateField(auto_now=True)),
                ('central', models.ForeignKey(on_delete=django.db.models.deletion.DO_NOTHING, to='core.central')),
            ],
            options={
                'db_table': 'versao_online',
            },
        ),
    ]
