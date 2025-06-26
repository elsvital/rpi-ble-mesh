# This is an auto-generated Django model module.
# You'll have to do the following manually to clean this up:
#   * Rearrange models' order
#   * Make sure each model has one field with primary_key=True
#   * Make sure each ForeignKey and OneToOneField has `on_delete` set to the desired behavior
#   * Remove `managed = False` lines if you wish to allow Django to create, modify, and delete the table
# Feel free to rename the models, but don't rename db_table values or field names.
from django.db import models
from django.utils import timezone

class AuthGroup(models.Model):
    name = models.CharField(unique=True, max_length=150)

    class Meta:
        managed = False
        db_table = 'auth_group'


class AuthGroupPermissions(models.Model):
    group = models.ForeignKey(AuthGroup, models.DO_NOTHING)
    permission = models.ForeignKey('AuthPermission', models.DO_NOTHING)

    class Meta:
        managed = False
        db_table = 'auth_group_permissions'
        unique_together = (('group', 'permission'),)


class AuthPermission(models.Model):
    content_type = models.ForeignKey('DjangoContentType', models.DO_NOTHING)
    codename = models.CharField(max_length=100)
    name = models.CharField(max_length=255)

    class Meta:
        managed = False
        db_table = 'auth_permission'
        unique_together = (('content_type', 'codename'),)


class AuthUser(models.Model):
    password = models.CharField(max_length=128)
    last_login = models.DateTimeField(blank=True, null=True)
    is_superuser = models.BooleanField()
    username = models.CharField(unique=True, max_length=150)
    last_name = models.CharField(max_length=150)
    email = models.CharField(max_length=254)
    is_staff = models.BooleanField()
    is_active = models.BooleanField()
    date_joined = models.DateTimeField()
    first_name = models.CharField(max_length=150)

    class Meta:
        managed = False
        db_table = 'auth_user'


class AuthUserGroups(models.Model):
    user = models.ForeignKey(AuthUser, models.DO_NOTHING)
    group = models.ForeignKey(AuthGroup, models.DO_NOTHING)

    class Meta:
        managed = False
        db_table = 'auth_user_groups'
        unique_together = (('user', 'group'),)


class AuthUserUserPermissions(models.Model):
    user = models.ForeignKey(AuthUser, models.DO_NOTHING)
    permission = models.ForeignKey(AuthPermission, models.DO_NOTHING)

    class Meta:
        managed = False
        db_table = 'auth_user_user_permissions'
        unique_together = (('user', 'permission'),)


class DjangoAdminLog(models.Model):
    object_id = models.TextField(blank=True, null=True)
    object_repr = models.CharField(max_length=200)
    action_flag = models.PositiveSmallIntegerField()
    change_message = models.TextField()
    content_type = models.ForeignKey('DjangoContentType', models.DO_NOTHING, blank=True, null=True)
    user = models.ForeignKey(AuthUser, models.DO_NOTHING)
    action_time = models.DateTimeField()

    class Meta:
        managed = False
        db_table = 'django_admin_log'


class DjangoContentType(models.Model):
    app_label = models.CharField(max_length=100)
    model = models.CharField(max_length=100)

    class Meta:
        managed = False
        db_table = 'django_content_type'
        unique_together = (('app_label', 'model'),)


class DjangoMigrations(models.Model):
    app = models.CharField(max_length=255)
    name = models.CharField(max_length=255)
    applied = models.DateTimeField()

    class Meta:
        managed = False
        db_table = 'django_migrations'


class DjangoSession(models.Model):
    session_key = models.CharField(primary_key=True, max_length=40)
    session_data = models.TextField()
    expire_date = models.DateTimeField()

    class Meta:
        managed = False
        db_table = 'django_session'

class Central(models.Model):
    id = models.CharField(max_length=254, primary_key=True)
    nome = models.CharField(max_length=254)
    versao = models.CharField(max_length=254)
    data_atualizacao = models.DateField(auto_now=True)

    class Meta:
        db_table = 'central'

class HistoricoAtualizacao(models.Model):
    micro = models.ForeignKey('MicroControlador', models.DO_NOTHING)
    data_atualizacao = models.DateField(blank=True, null=True)
    versao = models.CharField(max_length=254)
    log = models.CharField(max_length=254)

    class Meta:
        db_table = 'historico_atualizacao'


class MicroControlador(models.Model):
    id = models.CharField(max_length=254, primary_key=True)
    tipo_controlador = models.CharField(max_length=254)
    caminho_binario = models.CharField(max_length=254)
    token = models.TextField(blank=False, null=False)
    ip = models.CharField(max_length=254, blank=True, null=True)
    data_atualizacao = models.DateField(auto_now=True)

    class Meta:
        db_table = 'micro_controlador'
    def __str__(self):
        return f'{self.id} - {self.tipo_controlador} - {self.caminho_binario}'

class CentralMicro(models.Model):
    versao = models.CharField(max_length=254)
    data_atualizacao = models.DateField(auto_now=True)
    central = models.ForeignKey('Central', models.DO_NOTHING)
    micro = models.ForeignKey('MicroControlador', models.DO_NOTHING)

    class Meta:
        db_table = 'central_micro'

    def __str__(self):
        return f'{self.central.id} - {self.micro.id}'



class Topicos(models.Model):
    nome = models.CharField(max_length=254)
    mensagem = models.TextField(blank=True, null=True)
    central = models.ForeignKey(Central, models.DO_NOTHING)
    data_atualizacao = models.DateField(auto_now=True)

    class Meta:
        db_table = 'topicos'


class Treino(models.Model):
    micro = models.ForeignKey(MicroControlador, models.DO_NOTHING)
    user_id = models.CharField(max_length=254)
    plan_id = models.CharField(max_length=254,blank=True, null=True)
    main_muscles = models.CharField(max_length=254, blank=True, null=True)
    exercise_id = models.CharField(max_length=254, blank=True, null=True)
    gym_id = models.CharField(max_length=254, blank=True, null=True)
    summary = models.CharField(max_length=254, blank=True, null=True)
    workout_score = models.FloatField(blank=True, null=True)
    exercise_score = models.FloatField(blank=True, null=True)
    resting_time = models.FloatField(blank=True, null=True)
    used_load = models.FloatField(blank=True, null=True)
    failed_reps = models.IntegerField(blank=True, null=True)
    total_reps = models.IntegerField(blank=True, null=True)
    total_series = models.IntegerField(blank=True, null=True)
    amplitude_score = models.FloatField(blank=True, null=True)
    detailed_amplitude_score = models.FloatField(blank=True, null=True)
    time_score = models.FloatField(blank=True, null=True)
    detailed_time_score = models.FloatField(blank=True, null=True)
    data_atualizacao = models.DateField(auto_now=True)

    class Meta:
        db_table = 'treino'

    def __str__(self):
        data_str = self.data_atualizacao.strftime("%Y-%m-%d %H:%M:%S") if self.data_atualizacao else "sem data"
        return f'Treino do usuário {self.user_id}, enviado do micro controlador {self.micro.id} em: {data_str}'

class VersaoOnline(models.Model):
    central = models.ForeignKey(Central, models.DO_NOTHING)
    caminho_binario = models.CharField(max_length=254)
    versao = models.CharField(max_length=254)
    data_versao = models.DateField(blank=True, null=True)
    tipo_controlador = models.CharField(max_length=254)
    data_atualizacao = models.DateField(auto_now=True)

    class Meta:
        db_table = 'versao_online'
    def __str__(self):
        return f'{self.central.id} -  {self.versao} - {self.tipo_controlador}'

class SessaoCliente(models.Model):
    user_id = models.CharField(max_length=254)
    jwt_token = models.TextField()
    created_at = models.DateTimeField(default=timezone.now)
    last_seen = models.DateTimeField(default=timezone.now)

    class Meta:
        db_table = 'sessao_cliente'
        verbose_name = 'Sessão de Cliente'
        verbose_name_plural = 'Sessões de Clientes'

    def __str__(self):
        return f"Sessão de {self.user_id} iniciada em {self.created_at}, ultima atualizacao em {self.last_seen}"